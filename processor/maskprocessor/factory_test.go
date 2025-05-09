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

package maskprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	require.Equal(t, componentType, factory.Type())

	expectedCfg := &Config{}
	cfg, ok := factory.CreateDefaultConfig().(*Config)
	require.True(t, ok)
	require.Equal(t, expectedCfg, cfg)
}

func TestCreateLogsProcessorNilConfig(t *testing.T) {
	_, err := createLogsProcessor(context.Background(), processortest.NewNopSettings(componentType), nil, consumertest.NewNop())
	require.Error(t, err)
}

func TestCreateMetricsProcessorNilConfig(t *testing.T) {
	_, err := createMetricsProcessor(context.Background(), processortest.NewNopSettings(componentType), nil, consumertest.NewNop())
	require.Error(t, err)
}

func TestCreateTracesProcessorNilConfig(t *testing.T) {
	_, err := createTracesProcessor(context.Background(), processortest.NewNopSettings(componentType), nil, consumertest.NewNop())
	require.Error(t, err)
}

func TestCreateLogsProcessor(t *testing.T) {
	cfg := createDefaultConfig()
	p, err := createLogsProcessor(context.Background(), processortest.NewNopSettings(componentType), cfg, consumertest.NewNop())
	require.NotNil(t, p)
	require.NoError(t, err)
}

func TestCreateMetricsProcessor(t *testing.T) {
	cfg := createDefaultConfig()
	p, err := createMetricsProcessor(context.Background(), processortest.NewNopSettings(componentType), cfg, consumertest.NewNop())
	require.NotNil(t, p)
	require.NoError(t, err)
}

func TestCreateTracesProcessor(t *testing.T) {
	cfg := createDefaultConfig()
	p, err := createTracesProcessor(context.Background(), processortest.NewNopSettings(componentType), cfg, consumertest.NewNop())
	require.NotNil(t, p)
	require.NoError(t, err)
}
