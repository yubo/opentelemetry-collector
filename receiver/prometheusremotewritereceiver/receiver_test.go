package prometheusremotewritereceiver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/model/pdata"
)

func TestPrometheusRemoteWriteReceiver_New(t *testing.T) {
	defaultConfig := createDefaultConfig().(*Config)
	type args struct {
		config       Config
		nextConsumer consumer.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "nil_nextConsumer",
			args: args{
				config: *defaultConfig,
			},
			wantErr: componenterror.ErrNilNextConsumer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(componenttest.NewNopReceiverCreateSettings(), tt.args.config, tt.args.nextConsumer)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPrometheusRemoteWriteReceiver_Flush(t *testing.T) {
	ctx := context.Background()
	cfg := createDefaultConfig().(*Config)
	nextConsumer := consumertest.NewNop()
	rcv, err := New(componenttest.NewNopReceiverCreateSettings(), *cfg, nextConsumer)
	assert.NoError(t, err)
	r := rcv.(*prometheusRemoteWriteReceiver)
	var metrics = pdata.NewMetrics()
	assert.Nil(t, r.Flush(ctx, metrics, nextConsumer))
	r.Start(ctx, componenttest.NewNopHost())
	r.Shutdown(ctx)
}

func TestPrometheusRemoteWriteReceiver_EndToEnd(t *testing.T) {
	addr := GetAvailableLocalAddress(t)

	cfg := &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Endpoint:         addr,
	}
	sink := new(consumertest.MetricsSink)
	rcv, err := New(componenttest.NewNopReceiverCreateSettings(), *cfg, sink)
	require.NoError(t, err)
	r := rcv.(*prometheusRemoteWriteReceiver)

	require.NoError(t, r.Start(context.Background(), componenttest.NewNopHost()))
	defer r.Shutdown(context.Background())

	promReq := GeneratePromWriteRequest()
	promReqBody := GeneratePromWriteRequestBody(t, promReq)
	req, err := http.NewRequest(
		PromWriteHTTPMethod,
		fmt.Sprintf("http://%s/%s", r.config.Endpoint, r.config.WriteURL),
		promReqBody)
	require.NoError(t, err)

	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)
	mdd := sink.AllMetrics()
	require.Len(t, mdd, 1)
	require.Equal(t, 1, mdd[0].ResourceMetrics().Len())
	require.Equal(t, 1, mdd[0].ResourceMetrics().At(0).InstrumentationLibraryMetrics().Len())
	require.Equal(t, 2, mdd[0].ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().Len())

	metric := mdd[0].ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().At(0)
	assert.Equal(t, "first", metric.Name())
	assert.Equal(t, pdata.MetricDataTypeGauge, metric.DataType())
	require.Equal(t, 2, metric.Gauge().DataPoints().Len())

	metric = mdd[0].ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().At(1)
	assert.Equal(t, "second", metric.Name())
	assert.Equal(t, pdata.MetricDataTypeGauge, metric.DataType())
	require.Equal(t, 2, metric.Gauge().DataPoints().Len())

}

// GeneratePromWriteRequest generates a Prometheus remote
// write request.
func GeneratePromWriteRequest() *prompb.WriteRequest {
	req := &prompb.WriteRequest{
		Timeseries: []prompb.TimeSeries{{
			Labels: []prompb.Label{
				{Name: model.MetricNameLabel, Value: "first"},
				{Name: "foo", Value: "bar"},
				{Name: "biz", Value: "baz"},
			},
			Samples: []prompb.Sample{
				{Value: 1.0, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)},
				{Value: 2.0, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)},
			},
		}, {
			Labels: []prompb.Label{
				{Name: model.MetricNameLabel, Value: "second"},
				{Name: "foo", Value: "qux"},
				{Name: "bar", Value: "baz"},
			},
			Samples: []prompb.Sample{
				{Value: 3.0, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)},
				{Value: 4.0, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)},
			},
		}},
	}
	return req
}

// GeneratePromWriteRequestBody generates a Prometheus remote
// write request body.
func GeneratePromWriteRequestBody(
	t require.TestingT,
	req *prompb.WriteRequest,
) io.Reader {
	return bytes.NewReader(GeneratePromWriteRequestBodyBytes(t, req))
}

// GeneratePromWriteRequestBodyBytes generates a Prometheus remote
// write request body.
func GeneratePromWriteRequestBodyBytes(
	t require.TestingT,
	req *prompb.WriteRequest,
) []byte {
	data, err := proto.Marshal(req)
	require.NoError(t, err)

	compressed := snappy.Encode(nil, data)
	return compressed
}

// GetAvailableLocalAddress finds an available local port and returns an endpoint
// describing it. The port is available for opening when this function returns
// provided that there is no race by some other code to grab the same port
// immediately.
func GetAvailableLocalAddress(t *testing.T) string {
	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err, "Failed to get a free local port")
	// There is a possible race if something else takes this same port before
	// the test uses it, however, that is unlikely in practice.
	defer ln.Close()
	return ln.Addr().String()
}
