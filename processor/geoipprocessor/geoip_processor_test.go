package geoipprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/model/pdata"

	"github.com/yubo/opentelemetry-collector/internal/coreinternal/testdata"
)

var (
	cfg = &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Field:             []string{"ip"},
		DatabaseFile:      "./testdata/GeoLite2-City.mmdb",
		TargetField:       "geoip",
		IgnoreMissing:     false,
		FirstOnly:         true,
		Properties: []string{
			"continent_name",
			"country_iso_code",
			"country_name",
			"region_iso_code",
			"region_name",
			"city_name",
			"location",
		},
	}
)

func TestGeoipProcessorAttributesUpsert(t *testing.T) {
	tests := []struct {
		name             string
		config           *Config
		sourceAttributes map[string]pdata.Value
		wantAttributes   map[string]pdata.Value
	}{
		{
			name:             "config_with_attributes_applied_on_nil",
			config:           cfg,
			sourceAttributes: nil,
			wantAttributes:   nil,
		},
		{
			name:             "config_with_attributes_applied_on_empty_resource",
			config:           cfg,
			sourceAttributes: map[string]pdata.Value{},
			wantAttributes:   map[string]pdata.Value{},
		},
		{
			name:   "config_attributes_applied_on_existing_geoip",
			config: cfg,
			sourceAttributes: map[string]pdata.Value{
				"ip": pdata.NewValueString("1.1.1.1"),
			},
			wantAttributes: map[string]pdata.Value{
				"ip":    pdata.NewValueString("1.1.1.1"),
				"geoip": pdata.NewValueString("1.1.1.1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory()
			{
				tp, err := factory.CreateTracesProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), tt.config, consumertest.NewNop())
				require.NoError(t, err)
				require.NotNil(t, tp)

				td := generateTraceData(tt.sourceAttributes)
				assert.NoError(t, tp.ConsumeTraces(context.Background(), td))

				// Ensure that the modified `td` has the attributes sorted:
				sortAttributes(td)
				require.Equal(t, generateTraceData(tt.wantAttributes), td)
			}

			// Test metrics consumer
			{
				mp, err := factory.CreateMetricsProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), cfg, consumertest.NewNop())
				require.NoError(t, err)
				require.NotNil(t, mp)
				md := generateMetricData(tt.sourceAttributes)
				assert.NoError(t, mp.ConsumeMetrics(context.Background(), md))
				// Ensure that the modified `md` has the attributes sorted:
				sortMetricAttributes(md)
				require.Equal(t, generateMetricData(tt.wantAttributes), md)
			}

			// Test logs consumer
			{
				tp, err := factory.CreateLogsProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), cfg, consumertest.NewNop())
				require.NoError(t, err)
				require.NotNil(t, tp)
				ld := generateLogData(tt.sourceAttributes)
				assert.NoError(t, tp.ConsumeLogs(context.Background(), ld))
				// Ensure that the modified `ld` has the attributes sorted:
				sortLogAttributes(ld)
				require.Equal(t, generateLogData(tt.wantAttributes), ld)
			}
		})
	}
}

func generateTraceData(attrs map[string]pdata.Value) pdata.Traces {
	td := testdata.GenerateTracesOneSpan()
	if attrs == nil {
		return td
	}
	span := td.ResourceSpans().At(0).InstrumentationLibrarySpans().At(0).Spans().At(0)
	pdata.NewAttributeMapFromMap(attrs).CopyTo(span.Attributes())
	span.Attributes().Sort()
	return td
}

func sortAttributes(td pdata.Traces) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		rs.Resource().Attributes().Sort()
		ilss := rs.InstrumentationLibrarySpans()
		for j := 0; j < ilss.Len(); j++ {
			spans := ilss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				spans.At(k).Attributes().Sort()
			}
		}
	}
}
func generateMetricData(attrs map[string]pdata.Value) pdata.Metrics {
	md := testdata.GenerateMetricsOneMetric()
	if attrs == nil {
		return md
	}
	dp0 := md.ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0)
	pdata.NewAttributeMapFromMap(attrs).CopyTo(dp0.Attributes())
	dp0.Attributes().Sort()
	return md
}

func sortMetricAttributes(md pdata.Metrics) {
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rs := rms.At(i)
		rs.Resource().Attributes().Sort()
		ilms := rs.InstrumentationLibraryMetrics()
		for j := 0; j < ilms.Len(); j++ {
			metrics := ilms.At(j).Metrics()
			for k := 0; k < metrics.Len(); k++ {
				m := metrics.At(k)

				switch m.DataType() {
				case pdata.MetricDataTypeGauge:
					dps := m.Gauge().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).Attributes().Sort()
					}
				case pdata.MetricDataTypeSum:
					dps := m.Sum().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).Attributes().Sort()
					}
				case pdata.MetricDataTypeHistogram:
					dps := m.Histogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).Attributes().Sort()
					}
				case pdata.MetricDataTypeExponentialHistogram:
					dps := m.ExponentialHistogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).Attributes().Sort()
					}
				case pdata.MetricDataTypeSummary:
					dps := m.Summary().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).Attributes().Sort()
					}
				}
			}
		}
	}
}

func generateLogData(attrs map[string]pdata.Value) pdata.Logs {
	ld := testdata.GenerateLogsOneLogRecord()
	if attrs == nil {
		return ld
	}
	logs := ld.ResourceLogs().At(0).InstrumentationLibraryLogs().At(0).LogRecords().At(0)
	pdata.NewAttributeMapFromMap(attrs).CopyTo(logs.Attributes())
	logs.Attributes().Sort()
	return ld
}

func sortLogAttributes(ld pdata.Logs) {
	rss := ld.ResourceLogs()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		rs.Resource().Attributes().Sort()
		ilss := rs.InstrumentationLibraryLogs()
		for j := 0; j < ilss.Len(); j++ {
			logs := ilss.At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				s := logs.At(k)
				s.Attributes().Sort()
			}
		}
	}
}
