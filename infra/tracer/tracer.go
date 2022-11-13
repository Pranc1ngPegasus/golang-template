package tracer

import (
	"fmt"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/Pranc1ngPegasus/golang-template/domain/configuration"
	"github.com/Pranc1ngPegasus/golang-template/domain/logger"
	domain "github.com/Pranc1ngPegasus/golang-template/domain/tracer"
	"github.com/google/wire"
	"go.opencensus.io/trace"
)

var _ domain.Tracer = (*Tracer)(nil)

var NewTracerSet = wire.NewSet(
	wire.Bind(new(domain.Tracer), new(*Tracer)),
	NewTracer,
)

type Tracer struct {
	exporter *stackdriver.Exporter
}

func NewTracer(
	logger logger.Logger,
	config configuration.Configuration,
) (*Tracer, error) {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: config.Config().GCPProjectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exporter: %w", err)
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	return &Tracer{
		exporter: exporter,
	}, nil
}

func (t *Tracer) Start() error {
	if err := t.exporter.StartMetricsExporter(); err != nil {
		return fmt.Errorf("failed to start metrics exporter: %w", err)
	}

	return nil
}

func (t *Tracer) Stop() error {
	defer func() {
		t.exporter.Flush()
		t.exporter.StopMetricsExporter()
	}()

	return nil
}
