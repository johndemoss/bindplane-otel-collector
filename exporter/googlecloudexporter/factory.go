// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package googlecloudexporter provides a wrapper around the official googlecloudexporter component that exposes some quality of life improvements in configuration
package googlecloudexporter

import (
	"context"
	"fmt"

	gcp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/zap"
)

// gcpFactory is the factory used to create the underlying gcp exporter
var gcpFactory = gcp.NewFactory()

// componentType is the type of the google cloud exporter
var componentType = component.MustNewType("googlecloud")
var batchProcessorType = component.MustNewType("batch")

const (
	// The stability level of the exporter. Matches the current exporter in contrib
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for the googlecloud exporter
func NewFactory(collectorVersion string) exporter.Factory {
	return exporter.NewFactory(
		componentType,
		createDefaultConfig(collectorVersion),
		exporter.WithMetrics(createMetricsExporter, stability),
		exporter.WithLogs(createLogsExporter, stability),
		exporter.WithTraces(createTracesExporter, stability),
	)
}

// createMetricsExporter creates a metrics exporter based on this config.
func createMetricsExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Metrics, error) {
	exporterConfig := cfg.(*Config)
	exporterConfig.setClientOptions()

	if err := exporterConfig.setProject(); err != nil {
		set.Logger.Error("Failed to set project automatically", zap.Error(err))
	}

	gcpExporter, err := gcpFactory.CreateMetrics(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	processors := []processor.Metrics{}
	processorConfigs := []component.Config{
		exporterConfig.BatchConfig,
	}

	processorFactories := []processor.Factory{
		batchprocessor.NewFactory(),
	}

	processorSettings := processor.Settings{
		ID:                component.NewIDWithName(batchProcessorType, "googlecloud/metrics/batch"),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Metrics = gcpExporter
	for i, processorConfig := range processorConfigs {
		factory := processorFactories[i]
		processor, err := factory.CreateMetrics(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create metrics processor %s: %w", set.ID.String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &googlecloudExporter{
		appendHost:        exporterConfig.AppendHost,
		metricsProcessors: processors,
		metricsExporter:   gcpExporter,
		metricsConsumer:   consumer,
	}, nil
}

// createLogExporter creates a logs exporter based on this config.
func createLogsExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Logs, error) {
	exporterConfig := cfg.(*Config)
	exporterConfig.setClientOptions()

	if err := exporterConfig.setProject(); err != nil {
		set.Logger.Error("Failed to set project automatically", zap.Error(err))
	}

	gcpExporter, err := gcpFactory.CreateLogs(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs exporter: %w", err)
	}

	processors := []processor.Logs{}
	processorConfigs := []component.Config{
		exporterConfig.BatchConfig,
	}

	processorFactories := []processor.Factory{
		batchprocessor.NewFactory(),
	}

	processorSettings := processor.Settings{
		ID:                component.NewIDWithName(batchProcessorType, "googlecloud/logs/batch"),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Logs = gcpExporter
	for i, processorConfig := range processorConfigs {
		factory := processorFactories[i]
		processor, err := factory.CreateLogs(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create logs processor %s: %w", set.ID.String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &googlecloudExporter{
		appendHost:     exporterConfig.AppendHost,
		logsProcessors: processors,
		logsExporter:   gcpExporter,
		logsConsumer:   consumer,
	}, nil
}

// createTracesExporter creates a traces exporter based on this config.
func createTracesExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	exporterConfig := cfg.(*Config)
	exporterConfig.setClientOptions()

	if err := exporterConfig.setProject(); err != nil {
		set.Logger.Error("Failed to set project automatically", zap.Error(err))
	}

	gcpExporter, err := gcpFactory.CreateTraces(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create traces exporter: %w", err)
	}

	processors := []processor.Traces{}
	processorConfigs := []component.Config{
		exporterConfig.BatchConfig,
	}

	processorFactories := []processor.Factory{
		batchprocessor.NewFactory(),
	}

	processorSettings := processor.Settings{
		ID:                component.NewIDWithName(batchProcessorType, "googlecloud/traces/batch"),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Traces = gcpExporter
	for i, processorConfig := range processorConfigs {
		factory := processorFactories[i]
		processor, err := factory.CreateTraces(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create traces processor %s: %w", set.ID.String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &googlecloudExporter{
		appendHost:       exporterConfig.AppendHost,
		tracesProcessors: processors,
		tracesExporter:   gcpExporter,
		tracesConsumer:   consumer,
	}, nil
}
