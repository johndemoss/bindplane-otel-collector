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

//go:generate mdatagen metadata.yaml

// Package awss3eventextension implements an extension that consumes S3 event notifications
// from SQS and processes the objects containing OTLP data.
//
// The extension polls an SQS queue for S3 event notifications. When an object creation
// event is received, the extension downloads the S3 object and writes it so the configured
// directory.
package awss3eventextension // import "github.com/observiq/bindplane-otel-collector/extension/awss3eventextension"
