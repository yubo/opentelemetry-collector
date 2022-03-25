package geoipprocessor

import (
	"context"

	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"

	"github.com/yubo/opentelemetry-collector/internal/coreinternal/attraction"
)

type geoipProcessor struct {
	logger      *zap.Logger
	attrProc    *attraction.AttrProc
	fields      []string
	targetField string
}

func (p *geoipProcessor) processAttributes(ctx context.Context, attrs pdata.AttributeMap) bool {
	for _, field := range p.fields {
		if av, found := attrs.Get(field); found {
			// insert
			// TODO: use geoip.*
			attrs.Insert(p.targetField, av)
			return true
		}
	}

	return false
}

func (p *geoipProcessor) processTraces(ctx context.Context, td pdata.Traces) (pdata.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		ilss := rss.At(i).InstrumentationLibrarySpans()
		for j := 0; j < ilss.Len(); j++ {
			spans := ilss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				if done := p.processAttributes(ctx, spans.At(k).Attributes()); done {
					return td, nil
				}
			}
		}
	}
	return td, nil
}

func (p *geoipProcessor) processMetrics(ctx context.Context, md pdata.Metrics) (pdata.Metrics, error) {
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		ilms := rms.At(i).InstrumentationLibraryMetrics()
		for j := 0; j < ilms.Len(); j++ {
			metrics := ilms.At(j).Metrics()
			for k := 0; k < metrics.Len(); k++ {
				if done := p.processMetricAttributes(ctx, metrics.At(k)); done {
					return md, nil
				}
			}
		}
	}
	return md, nil
}

// Attributes are provided for each log and trace, but not at the metric level
// Need to process attributes for every data point within a metric.
func (p *geoipProcessor) processMetricAttributes(ctx context.Context, m pdata.Metric) bool {

	// This is a lot of repeated code, but since there is no single parent superclass
	// between metric data types, we can't use polymorphism.
	switch m.DataType() {
	case pdata.MetricDataTypeGauge:
		dps := m.Gauge().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			if done := p.processAttributes(ctx, dps.At(i).Attributes()); done {
				return true
			}
		}
	case pdata.MetricDataTypeSum:
		dps := m.Sum().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			if done := p.processAttributes(ctx, dps.At(i).Attributes()); done {
				return true
			}
		}
	case pdata.MetricDataTypeHistogram:
		dps := m.Histogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			if done := p.processAttributes(ctx, dps.At(i).Attributes()); done {
				return true
			}
		}
	case pdata.MetricDataTypeExponentialHistogram:
		dps := m.ExponentialHistogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			if done := p.processAttributes(ctx, dps.At(i).Attributes()); done {
				return true
			}
		}
	case pdata.MetricDataTypeSummary:
		dps := m.Summary().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			if done := p.processAttributes(ctx, dps.At(i).Attributes()); done {
				return true
			}
		}
	}

	return false
}

func (p *geoipProcessor) processLogs(ctx context.Context, ld pdata.Logs) (pdata.Logs, error) {
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		ilss := rls.At(i).InstrumentationLibraryLogs()
		for j := 0; j < ilss.Len(); j++ {
			logs := ilss.At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				if done := p.processAttributes(ctx, logs.At(k).Attributes()); done {
					return ld, nil
				}
			}
		}
	}
	return ld, nil
}
