// Copyright The OpenTelemetry Authors
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

package ast // import "go.opentelemetry.io/otel/schema/v1.1/ast"

import (
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	types10 "go.opentelemetry.io/otel/schema/v1.0/types"
	types11 "go.opentelemetry.io/otel/schema/v1.1/types"
)

// Metrics corresponds to a section representing a list of changes that happened
// to metrics schema in a particular version.
type Metrics struct {
	Changes []MetricsChange
}

// MetricsChange corresponds to a section representing metrics change.
type MetricsChange struct {
	RenameMetrics    map[types10.MetricName]types10.MetricName `yaml:"rename_metrics"`
	RenameAttributes *ast10.AttributeMapForMetrics             `yaml:"rename_attributes"`
	Split            *SplitMetric                              `yaml:"split"`
}

// SplitMetric  corresponds to a section representing a splitting of a metric
// into multiple metrics by eliminating an attribute.
// SplitMetrics is introduced in schema file format 1.1,
// see https://github.com/open-telemetry/opentelemetry-specification/pull/2653
type SplitMetric struct {
	// Name of the old metric to split.
	ApplyToMetric types10.MetricName `yaml:"apply_to_metric"`

	// Name of attribute in the old metric to use for splitting. The attribute will be
	// eliminated, the new metric will not have it.
	ByAttribute types11.AttributeName `yaml:"by_attribute"`

	// Names of new metrics to create, one for each possible value of attribute.
	// map of key/values. The keys are the new metric name starting from this version,
	// the values are old attribute value used in the previous version.
	MetricsFromAttributes map[types10.MetricName]types11.AttributeValue `yaml:"metrics_from_attributes"`
}
