// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topologyprocessor

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/ptracetest"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

func TestProcessor_Logs(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")

	tmp, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	logs, err := golden.ReadLogs(filepath.Join("testdata", "logs", "w3c-logs.yaml"))
	require.NoError(t, err)

	ctx := client.NewContext(context.Background(), client.Info{
		Metadata: client.NewMetadata(map[string][]string{
			accountIDHeader:      {"myAccountID1"},
			organizationIDHeader: {"myOrgID1"},
			configurationHeader:  {"myConfigName1"},
			resourceNameHeader:   {"myResourceName1"},
		}),
	})
	processedLogs, err := tmp.processLogs(ctx, logs)
	require.NoError(t, err)

	// Output logs should be the same as input logs (passthrough check)
	require.NoError(t, plogtest.CompareLogs(logs, processedLogs))

	// validate that upsert route was performed
	require.True(t, tmp.topology.GatewaySource.AccountID == "myAccountID")
	require.True(t, tmp.topology.GatewaySource.OrganizationID == "myOrgID")
	require.True(t, tmp.topology.GatewaySource.Configuration == "myConfigName")
	ci := GatewayInfo{
		Configuration:  "myConfigName1",
		AccountID:      "myAccountID1",
		OrganizationID: "myOrgID1",
		GatewayID:      "myResourceName1",
	}
	_, ok := tmp.topology.RouteTable[ci]
	require.True(t, ok)
}

func TestProcessor_Metrics(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")

	tmp, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	metrics, err := golden.ReadMetrics(filepath.Join("testdata", "metrics", "host-metrics.yaml"))
	require.NoError(t, err)

	ctx := client.NewContext(context.Background(), client.Info{
		Metadata: client.NewMetadata(map[string][]string{
			accountIDHeader:      {"myAccountID1"},
			organizationIDHeader: {"myOrgID1"},
			configurationHeader:  {"myConfigName1"},
			resourceNameHeader:   {"myResourceName1"},
		}),
	})

	processedMetrics, err := tmp.processMetrics(ctx, metrics)
	require.NoError(t, err)

	// Output metrics should be the same as input logs (passthrough check)
	require.NoError(t, pmetrictest.CompareMetrics(metrics, processedMetrics))

	// validate that upsert route was performed
	require.True(t, tmp.topology.GatewaySource.AccountID == "myAccountID")
	require.True(t, tmp.topology.GatewaySource.OrganizationID == "myOrgID")
	require.True(t, tmp.topology.GatewaySource.Configuration == "myConfigName")
	ci := GatewayInfo{
		Configuration:  "myConfigName1",
		AccountID:      "myAccountID1",
		OrganizationID: "myOrgID1",
		GatewayID:      "myResourceName1",
	}
	_, ok := tmp.topology.RouteTable[ci]
	require.True(t, ok)
}

func TestProcessor_Traces(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")

	tmp, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	traces, err := golden.ReadTraces(filepath.Join("testdata", "traces", "bindplane-traces.yaml"))
	require.NoError(t, err)

	ctx := client.NewContext(context.Background(), client.Info{
		Metadata: client.NewMetadata(map[string][]string{
			accountIDHeader:      {"myAccountID1"},
			organizationIDHeader: {"myOrgID1"},
			configurationHeader:  {"myConfigName1"},
			resourceNameHeader:   {"myResourceName1"},
		}),
	})

	processedTraces, err := tmp.processTraces(ctx, traces)
	require.NoError(t, err)

	// Output traces should be the same as input logs (passthrough check)
	require.NoError(t, ptracetest.CompareTraces(traces, processedTraces))

	// validate that upsert route was performed
	require.True(t, tmp.topology.GatewaySource.AccountID == "myAccountID")
	require.True(t, tmp.topology.GatewaySource.OrganizationID == "myOrgID")
	require.True(t, tmp.topology.GatewaySource.Configuration == "myConfigName")
	ci := GatewayInfo{
		Configuration:  "myConfigName1",
		AccountID:      "myAccountID1",
		OrganizationID: "myOrgID1",
		GatewayID:      "myResourceName1",
	}
	_, ok := tmp.topology.RouteTable[ci]
	require.True(t, ok)
}

func TestProcessor_MissingHeader(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")

	tmp, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	traces, err := golden.ReadTraces(filepath.Join("testdata", "traces", "bindplane-traces.yaml"))
	require.NoError(t, err)

	ctx := client.NewContext(context.Background(), client.Info{
		Metadata: client.NewMetadata(map[string][]string{
			organizationIDHeader: {"myOrgID1"},
			configurationHeader:  {"myConfigName1"},
			resourceNameHeader:   {"myResourceName1"},
		}),
	})

	processedTraces, err := tmp.processTraces(ctx, traces)
	require.NoError(t, err)

	// Output traces should be the same as input logs (passthrough check)
	require.NoError(t, ptracetest.CompareTraces(traces, processedTraces))

	// validate that upsert route was not performed
	require.Equal(t, 0, len(tmp.topology.RouteTable))
}

// Test 2 instances with the same processor ID
func TestProcessor_Logs_TwoInstancesSameID(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")

	tmp1, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	tmp2, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID2",
		AccountID:      "myAccountID2",
		Configuration:  "myConfigName2",
	}, processorID)
	require.NoError(t, err)

	logs, err := golden.ReadLogs(filepath.Join("testdata", "logs", "w3c-logs.yaml"))
	require.NoError(t, err)

	_, err = tmp1.processLogs(context.Background(), logs)
	require.NoError(t, err)

	_, err = tmp2.processLogs(context.Background(), logs)
	require.NoError(t, err)
}

func TestProcessor_Logs_TwoInstancesDifferentID(t *testing.T) {
	processorID := component.MustNewIDWithName("topology", "1")
	processorID2 := component.MustNewIDWithName("topology", "2")

	tmp1, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID",
		AccountID:      "myAccountID",
		Configuration:  "myConfigName",
	}, processorID)
	require.NoError(t, err)

	tmp2, err := newTopologyProcessor(zap.NewNop(), &Config{
		OrganizationID: "myOrgID2",
		AccountID:      "myAccountID2",
		Configuration:  "myConfigName2",
	}, processorID2)
	require.NoError(t, err)

	logs, err := golden.ReadLogs(filepath.Join("testdata", "logs", "w3c-logs.yaml"))
	require.NoError(t, err)

	_, err = tmp1.processLogs(context.Background(), logs)
	require.NoError(t, err)

	_, err = tmp2.processLogs(context.Background(), logs)
	require.NoError(t, err)
}
