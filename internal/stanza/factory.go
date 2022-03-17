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
	"context"

	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/pipeline"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
)

// LogReceiverType is the interface used by stanza-based log receivers
type LogReceiverType interface {
	Type() config.Type
	CreateDefaultConfig() config.Receiver
	BaseConfig(config.Receiver) BaseConfig
	DecodeInputConfig(config.Receiver) (*operator.Config, error)
}

// NewFactory creates a factory for a Stanza-based receiver
func NewFactory(logReceiverType LogReceiverType) component.ReceiverFactory {
	return component.NewReceiverFactory(
		logReceiverType.Type(),
		logReceiverType.CreateDefaultConfig,
		component.WithLogsReceiver(createLogsReceiver(logReceiverType)),
	)
}

func createLogsReceiver(logReceiverType LogReceiverType) component.CreateLogsReceiverFunc {
	return func(
		ctx context.Context,
		params component.ReceiverCreateSettings,
		cfg config.Receiver,
		nextConsumer consumer.Logs,
	) (component.LogsReceiver, error) {
		inputCfg, err := logReceiverType.DecodeInputConfig(cfg)
		if err != nil {
			return nil, err
		}

		baseCfg := logReceiverType.BaseConfig(cfg)
		operatorCfgs, err := baseCfg.decodeOperatorConfigs()
		if err != nil {
			return nil, err
		}

		operators := append([]operator.Config{*inputCfg}, operatorCfgs...)

		emitterOpts := []LogEmitterOption{
			LogEmitterWithLogger(params.Logger.Sugar()),
		}

		if baseCfg.Converter.MaxFlushCount > 0 {
			emitterOpts = append(emitterOpts, LogEmitterWithMaxBatchSize(baseCfg.Converter.MaxFlushCount))
		}

		if baseCfg.Converter.FlushInterval > 0 {
			emitterOpts = append(emitterOpts, LogEmitterWithFlushInterval(baseCfg.Converter.FlushInterval))
		}

		emitter := NewLogEmitter(emitterOpts...)
		pipe, err := pipeline.Config{
			Operators:     operators,
			DefaultOutput: emitter,
		}.Build(params.Logger.Sugar())
		if err != nil {
			return nil, err
		}

		opts := []ConverterOption{
			WithLogger(params.Logger),
		}

		if baseCfg.Converter.WorkerCount > 0 {
			opts = append(opts, WithWorkerCount(baseCfg.Converter.WorkerCount))
		}
		converter := NewConverter(opts...)
		obsrecv := obsreport.NewReceiver(obsreport.ReceiverSettings{
			ReceiverID:             cfg.ID(),
			ReceiverCreateSettings: params,
		})
		return &receiver{
			id:        cfg.ID(),
			pipe:      pipe,
			emitter:   emitter,
			consumer:  nextConsumer,
			logger:    params.Logger,
			converter: converter,
			obsrecv:   obsrecv,
		}, nil
	}
}
