package trace

import (
	"context"

	"github.com/google/uuid"
)

// context key for TraceID
type contextKey string

const traceIDKey contextKey = "traceID"

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func GetTraceID(ctx context.Context) string {
	val := ctx.Value(traceIDKey)
	if id, ok := val.(string); ok {
		return id
	}
	return ""
}

func NewTraceID() string {
	return "trace-" + uuid.New().String()
}
