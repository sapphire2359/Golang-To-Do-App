package api

import (
	"net/http"
	"todoapp-json/trace"
)

// middleware func that adds a TraceID to each request's context.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		ctxWithTrace := trace.WithTraceID(request.Context(), trace.NewTraceID())
		request = request.WithContext(ctxWithTrace)
		next.ServeHTTP(write, request)
	})
}
