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

// Package topologyprocessor collects metrics, traces, and logs for
package topologyprocessor

import (
	"errors"

	"go.opentelemetry.io/collector/component"
)

// Config is the configuration for the processor
type Config struct {
	// Bindplane extension to use in order to report topology. Optional.
	BindplaneExtension *component.ID `mapstructure:"bindplane_extension"`

	// Name of the Config where this processor is present
	Configuration string `mapstructure:"configuration"`

	// OrganizationID of the Org where this processor is present
	OrganizationID string `mapstructure:"organizationID"`

	// AccountID of the Account where this processor is present
	AccountID string `mapstructure:"accountID"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	if cfg.Configuration == "" {
		return errors.New("`configuration` must be specified")
	}

	if cfg.OrganizationID == "" {
		return errors.New("`organizationID` must be specified")
	}

	if cfg.AccountID == "" {
		return errors.New("`accountID` must be specified")
	}

	return nil
}
