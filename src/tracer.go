package src

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type contextIDKey struct{}

func NewTraceContext(ctx context.Context) context.Context {
	id := uuid.NewV1()
	return context.WithValue(ctx, contextIDKey{}, id)
}

func GetIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(contextIDKey{}).(uuid.UUID)
}
