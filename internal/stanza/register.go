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

package stanza // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza"

import (
	// Register parsers and transformers for stanza-based log receivers
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/csv"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/json"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/regex"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/severity"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/time"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/trace"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/parser/uri"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/add"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/copy"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/filter"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/flatten"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/metadata"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/move"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/recombine"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/restructure"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/retain"
	_ "github.com/open-telemetry/opentelemetry-log-collection/operator/transformer/router"
)
