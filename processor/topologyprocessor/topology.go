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
	"fmt"
	"sync"
	"time"
)

// TopoRegistry represents a registry for the topology processor to register their TopologyState.
type TopoRegistry interface {
	// RegisterTopologyState registers the topology state for the given processor.
	// It should return an error if the processor has already been registered.
	RegisterTopologyState(processorID string, data *TopoState) error
	Reset()
}

// TopoState represents the data captured through topology processors.
type TopoState struct {
	// GatewaySource is the gateway source that the entries in the route table point to
	GatewaySource GatewayInfo
	// RouteTable is a map of gateway destinations to the time at which they were last detected
	RouteTable map[GatewayInfo]time.Time
	mux        sync.RWMutex
}

// GatewayInfo represents a bindplane gateway source or destination
type GatewayInfo struct {
	// OrganizationID is the organizationID where this gateway dest/source lives
	OrganizationID string `json:"organizationID"`
	// AccountID is the accountID where this gateway dest/source lives
	AccountID string `json:"accountID"`
	// Configuration is the name of the configuration where this gateway dest/source lives
	Configuration string `json:"configuration"`
	// GatewayID is the ComponentID of a gateway source, or the resource name of a gateway destination
	GatewayID string `json:"gatewayID"`
}

// GatewayRecord represents a gateway destination and the time it was last detected
type GatewayRecord struct {
	// Gateway represents a gateway destinations
	Gateway GatewayInfo `json:"gateway"`
	// LastUpdated is a timestamp of the last time a message w/ topology headers was detected from the above gateway destination
	LastUpdated time.Time `json:"lastUpdated"`
}

// TopoInfo represents a gateway source & the gateway destinations that point to it.
type TopoInfo struct {
	GatewaySource       GatewayInfo     `json:"gatewaySource"`
	GatewayDestinations []GatewayRecord `json:"gatewayDestinations"`
}

// NewTopologyState initializes a new TopologyState
func NewTopologyState(gw GatewayInfo) (*TopoState, error) {
	return &TopoState{
		GatewaySource: gw,
		RouteTable:    make(map[GatewayInfo]time.Time),
		mux:           sync.RWMutex{},
	}, nil
}

// UpsertRoute upserts given route.
func (ts *TopoState) UpsertRoute(_ context.Context, gw GatewayInfo) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	ts.RouteTable[gw] = time.Now()
}

// ResettableTopologyRegistry is a concrete version of TopologyDataRegistry that is able to be reset.
type ResettableTopologyRegistry struct {
	topology *sync.Map
}

// NewResettableTopologyRegistry creates a new ResettableTopologyRegistry
func NewResettableTopologyRegistry() *ResettableTopologyRegistry {
	return &ResettableTopologyRegistry{
		topology: &sync.Map{},
	}
}

// RegisterTopologyState registers the TopologyState with the registry.
func (rtsr *ResettableTopologyRegistry) RegisterTopologyState(processorID string, topoState *TopoState) error {
	_, alreadyExists := rtsr.topology.LoadOrStore(processorID, topoState)
	if alreadyExists {
		return fmt.Errorf("topology for processor %q was already registered", processorID)
	}

	return nil
}

// Reset unregisters all topology states in this registry
func (rtsr *ResettableTopologyRegistry) Reset() {
	rtsr.topology = &sync.Map{}
}

// TopologyInfos returns all the topology data in this registry.
func (rtsr *ResettableTopologyRegistry) TopologyInfos() []TopoInfo {
	states := []*TopoState{}

	rtsr.topology.Range(func(_, value any) bool {
		ts := value.(*TopoState)
		states = append(states, ts)
		return true
	})

	ti := []TopoInfo{}
	for _, ts := range states {
		curInfo := getTopoInfoFromState(ts)

		if len(curInfo.GatewayDestinations) > 0 {
			ti = append(ti, curInfo)
		}
	}

	return ti
}

func getTopoInfoFromState(ts *TopoState) TopoInfo {
	ts.mux.RLock()
	defer ts.mux.RUnlock()
	curInfo := TopoInfo{}
	curInfo.GatewaySource.OrganizationID = ts.GatewaySource.OrganizationID
	curInfo.GatewaySource.AccountID = ts.GatewaySource.AccountID
	curInfo.GatewaySource.Configuration = ts.GatewaySource.Configuration
	curInfo.GatewaySource.GatewayID = ts.GatewaySource.GatewayID
	for gw, updated := range ts.RouteTable {
		curInfo.GatewayDestinations = append(curInfo.GatewayDestinations, GatewayRecord{
			Gateway: GatewayInfo{
				OrganizationID: gw.OrganizationID,
				AccountID:      gw.AccountID,
				Configuration:  gw.Configuration,
				GatewayID:      gw.GatewayID,
			},
			LastUpdated: updated.UTC(),
		})
	}

	return curInfo
}
