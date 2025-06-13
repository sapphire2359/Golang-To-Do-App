package trace

import (
	"context"
	"fmt"
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
	return fmt.Sprintf("trace-%d", RandInt()) // Replace with UUID if needed
}

func RandInt() int {
	return int(1000 + (9999-1000)*int64(1))
}
