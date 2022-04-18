package prometheusremotewritereceiver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/value"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
)

const (
	// PromWriteURL is the url for the prom write handler
	PromWriteURL = "/prom/remote/write"

	// PromWriteHTTPMethod is the HTTP method used with this resource.
	PromWriteHTTPMethod = http.MethodPost
)

var _ component.MetricsReceiver = (*prometheusRemoteWriteReceiver)(nil)
var pdataStaleFlags = pdata.NewMetricDataPointFlags(pdata.MetricDataPointFlagNoRecordedValue)

// prometheusRemoteWriteReceiver implements the component.MetricsReceiver for StatsD protocol.
type prometheusRemoteWriteReceiver struct {
	settings component.ReceiverCreateSettings
	config   *Config

	server       *http.Server
	consumer     consumer.Metrics
	nodeResource *pdata.Resource
	ctx          context.Context
	cancel       context.CancelFunc
}

// New creates the StatsD receiver with the given parameters.
func New(
	set component.ReceiverCreateSettings,
	config Config,
	next consumer.Metrics,
) (component.MetricsReceiver, error) {
	if next == nil {
		return nil, componenterror.ErrNilNextConsumer
	}

	r := &prometheusRemoteWriteReceiver{
		config:       &config,
		settings:     set,
		consumer:     next,
		nodeResource: nil, // TODO
	}
	return r, nil
}

// Start
func (p *prometheusRemoteWriteReceiver) Start(ctx context.Context, host component.Host) error {
	p.ctx, p.cancel = context.WithCancel(ctx)

	p.server = &http.Server{Addr: p.config.Endpoint, Handler: p}

	go func() {
		if err := p.server.ListenAndServe(); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				host.ReportFatalError(err)
			}
		}
	}()

	return nil
}

// Shutdown stops the StatsD receiver.
func (p *prometheusRemoteWriteReceiver) Shutdown(context.Context) error {
	err := p.server.Close()
	p.cancel()
	return err
}

func (p *prometheusRemoteWriteReceiver) Flush(ctx context.Context, metrics pdata.Metrics, nextConsumer consumer.Metrics) error {
	if err := nextConsumer.ConsumeMetrics(ctx, metrics); err != nil {
		return err
	}

	return nil
}

func (p *prometheusRemoteWriteReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := remote.DecodeWriteRequest(r.Body)
	if err != nil {
		p.settings.Logger.Error("Error decoding remote write request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = p.write(req)
	if err != nil {
		p.settings.Logger.Error("Error appending remote write", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *prometheusRemoteWriteReceiver) write(req *prompb.WriteRequest) error {

	if metrics := p.reqToMetrics(req); metrics != nil {
		return p.consumer.ConsumeMetrics(p.ctx, *metrics)
	}

	return nil
}
func (p *prometheusRemoteWriteReceiver) reqToMetrics(req *prompb.WriteRequest) *pdata.Metrics {
	metrics := pdata.NewMetricSlice()
	for _, ts := range req.Timeseries {
		metric := pdata.NewMetric()
		metric.SetDataType(pdata.MetricDataTypeGauge)

		for _, label := range ts.Labels {
			if label.Name == model.MetricNameLabel {
				metric.SetName(label.Value)
			}
		}

		//metric.SetDescription(mf.metadata.Help)
		//metric.SetUnit(mf.metadata.Unit)
		gdpL := metric.Gauge().DataPoints()

		for _, sample := range ts.Samples {
			if err := toNumberDataPoint(sample, ts.Labels, &gdpL); err != nil {
				p.settings.Logger.Error("toNumberDataPointSlice", zap.Error(err))
			}
		}
		metric.CopyTo(metrics.AppendEmpty())
	}

	return p.metricSliceToMetrics(&metrics)
	//t.sink.ConsumeMetrics(ctx, *metrics)
}

func (p *prometheusRemoteWriteReceiver) metricSliceToMetrics(metricsL *pdata.MetricSlice) *pdata.Metrics {
	if metricsL.Len() == 0 {
		return nil
	}

	metrics := pdata.NewMetrics()
	rms := metrics.ResourceMetrics().AppendEmpty()
	ilm := rms.InstrumentationLibraryMetrics().AppendEmpty()
	metricsL.CopyTo(ilm.Metrics())
	//p.nodeResource.CopyTo(rms.Resource())
	return &metrics
}

func toNumberDataPoint(sample prompb.Sample, labels []prompb.Label, dest *pdata.NumberDataPointSlice) error {
	tsNanos := pdataTimestampFromMs(sample.Timestamp)

	point := dest.AppendEmpty()
	point.SetTimestamp(tsNanos)

	if value.IsStaleNaN(sample.Value) {
		point.SetFlags(pdataStaleFlags)
	} else {
		point.SetDoubleVal(sample.Value)
	}

	// labels
	attrs := point.Attributes()
	for _, label := range labels {
		attrs.InsertString(label.Name, label.Value)
	}

	return nil
}

func pdataTimestampFromMs(timeAtMs int64) pdata.Timestamp {
	secs, ns := timeAtMs/1e3, (timeAtMs%1e3)*1e6
	return pdata.NewTimestampFromTime(time.Unix(secs, ns))
}
