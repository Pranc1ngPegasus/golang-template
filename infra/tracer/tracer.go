package tracer

import (
	"context"
	"fmt"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/Pranc1ngPegasus/golang-template/domain/configuration"
	"github.com/Pranc1ngPegasus/golang-template/domain/logger"
	domain "github.com/Pranc1ngPegasus/golang-template/domain/tracer"
	"github.com/google/wire"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/multierr"
)

var _ domain.Tracer = (*Tracer)(nil)

var NewTracerSet = wire.NewSet(
	wire.Bind(new(domain.Tracer), new(*Tracer)),
	NewTracer,
)

type Tracer struct {
	exporter *cloudtrace.Exporter
	provider *sdktrace.TracerProvider
}

func NewTracer(
	logger logger.Logger,
	config configuration.Configuration,
) (*Tracer, error) {
	cfg := config.Config()
	ctx := context.Background()

	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(cfg.GCPProjectID))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(
			gcp.NewDetector(),
		),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return &Tracer{
		exporter: exporter,
		provider: tp,
	}, nil
}

func (t *Tracer) Stop() (err error) {
	ctx := context.Background()

	defer func() {
		var merr error

		multierr.AppendInto(&merr, t.exporter.Shutdown(ctx))
		multierr.AppendInto(&merr, t.provider.ForceFlush(ctx))
		multierr.AppendInto(&merr, t.provider.Shutdown(ctx))

		err = merr
	}()

	return nil
}
