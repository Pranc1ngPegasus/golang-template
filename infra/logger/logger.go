package logger

import (
	"context"
	"fmt"

	domain "github.com/Pranc1ngPegasus/golang-template/domain/logger"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var _ domain.Logger = (*Logger)(nil)

type Logger struct {
	logger *otelzap.Logger
}

var NewLoggerSet = wire.NewSet(
	wire.Bind(new(domain.Logger), new(*Logger)),
	NewLogger,
)

func NewLogger() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger := otelzap.New(
		log,
		otelzap.WithTraceIDField(true),
	)

	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) Field(key string, message interface{}) domain.Field {
	return domain.Field{
		Key:       key,
		Interface: message,
	}
}

func (l *Logger) field(field domain.Field) zap.Field {
	switch i := field.Interface.(type) {
	case string:
		return zap.String(field.Key, i)
	case int:
		return zap.Int(field.Key, i)
	case bool:
		return zap.Bool(field.Key, i)
	default:
		return zap.Any(field.Key, i)
	}
}

func (l *Logger) Info(ctx context.Context, message string, fields ...domain.Field) {
	zapfields := lo.Map(fields, func(field domain.Field, _ int) zap.Field {
		return l.field(field)
	})

	l.logger.Ctx(ctx).Info(message, zapfields...)
}

func (l *Logger) Error(ctx context.Context, message string, err error, fields ...domain.Field) {
	zapfields := lo.Map(fields, func(field domain.Field, _ int) zap.Field {
		return l.field(field)
	})

	zapfields = append(zapfields, zap.Error(err))

	l.logger.Ctx(ctx).Error(message, zapfields...)
}
