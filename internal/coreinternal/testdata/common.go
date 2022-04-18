// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testdata

import (
	"go.opentelemetry.io/collector/model/pdata"
)

var (
	resourceAttributes1 = map[string]pdata.Value{"resource-attr": pdata.NewValueString("resource-attr-val-1")}
	resourceAttributes2 = map[string]pdata.Value{"resource-attr": pdata.NewValueString("resource-attr-val-2")}
	spanEventAttributes = map[string]pdata.Value{"span-event-attr": pdata.NewValueString("span-event-attr-val")}
	spanLinkAttributes  = map[string]pdata.Value{"span-link-attr": pdata.NewValueString("span-link-attr-val")}
	spanAttributes      = map[string]pdata.Value{"span-attr": pdata.NewValueString("span-attr-val")}
)

const (
	TestLabelKey1       = "label-1"
	TestLabelValue1     = "label-value-1"
	TestLabelKey2       = "label-2"
	TestLabelValue2     = "label-value-2"
	TestLabelKey3       = "label-3"
	TestLabelValue3     = "label-value-3"
	TestAttachmentKey   = "exemplar-attachment"
	TestAttachmentValue = "exemplar-attachment-value"
)

func initResourceAttributes1(dest pdata.AttributeMap) {
	pdata.NewAttributeMapFromMap(resourceAttributes1).CopyTo(dest)
}

func initResourceAttributes2(dest pdata.AttributeMap) {
	pdata.NewAttributeMapFromMap(resourceAttributes2).CopyTo(dest)
}

func initSpanAttributes(dest pdata.AttributeMap) {
	pdata.NewAttributeMapFromMap(spanAttributes).CopyTo(dest)
}

func initSpanEventAttributes(dest pdata.AttributeMap) {
	pdata.NewAttributeMapFromMap(spanEventAttributes).CopyTo(dest)
}

func initSpanLinkAttributes(dest pdata.AttributeMap) {
	pdata.NewAttributeMapFromMap(spanLinkAttributes).CopyTo(dest)
}

func initMetricAttachment(dest pdata.AttributeMap) {
	dest.UpsertString(TestAttachmentKey, TestAttachmentValue)
}

func initMetricAttributes1(dest pdata.AttributeMap) {
	dest.UpsertString(TestLabelKey1, TestLabelValue1)
}

func initMetricAttributes12(dest pdata.AttributeMap) {
	dest.UpsertString(TestLabelKey1, TestLabelValue1)
	dest.UpsertString(TestLabelKey2, TestLabelValue2)
	dest.Sort()
}

func initMetricAttributes13(dest pdata.AttributeMap) {
	dest.UpsertString(TestLabelKey1, TestLabelValue1)
	dest.UpsertString(TestLabelKey3, TestLabelValue3)
	dest.Sort()
}

func initMetricAttributes2(dest pdata.AttributeMap) {
	dest.UpsertString(TestLabelKey2, TestLabelValue2)
}
