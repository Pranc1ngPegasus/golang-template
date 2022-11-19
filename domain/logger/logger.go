//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

package logger

import (
	"context"
)

type Logger interface {
	Field(string, interface{}) Field
	Info(context.Context, string, ...Field)
	Error(context.Context, string, error, ...Field)
}

type (
	Field struct {
		Key       string
		Interface interface{}
	}
)
