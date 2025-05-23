// Copyright observIQ, Inc.
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

package chronicleforwarderexporter

import (
	"context"
	"errors"
	"time"

	"github.com/observiq/bindplane-otel-collector/exporter/chronicleforwarderexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// NewFactory creates a new Chronicle exporter factory.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability))
}

// createDefaultConfig creates the default configuration for the exporter.
func createDefaultConfig() component.Config {
	return &Config{
		TimeoutConfig:    exporterhelper.NewDefaultTimeoutConfig(),
		QueueBatchConfig: exporterhelper.NewDefaultQueueConfig(),
		BackOffConfig:    configretry.NewDefaultBackOffConfig(),
		ExportType:       exportTypeSyslog,
		Syslog: SyslogConfig{
			AddrConfig: confignet.AddrConfig{
				Endpoint:  "127.0.0.1:10514",
				Transport: "tcp",
				DialerConfig: confignet.DialerConfig{
					Timeout: 5 * time.Second,
				},
			},
		},
	}
}

// createLogsExporter creates a new log exporter based on this config.
func createLogsExporter(
	ctx context.Context,
	params exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	forwarderCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	exp, err := newExporter(forwarderCfg, params)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewLogs(
		ctx,
		params,
		forwarderCfg,
		exp.logsDataPusher,
		exporterhelper.WithCapabilities(exp.Capabilities()),
		exporterhelper.WithTimeout(forwarderCfg.TimeoutConfig),
		exporterhelper.WithQueue(forwarderCfg.QueueBatchConfig),
		exporterhelper.WithRetry(forwarderCfg.BackOffConfig),
	)
}
