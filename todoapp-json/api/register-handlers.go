package api

import (
	"net/http"
)

// registering all handlers wrapped with middleware trace function
func RegisterHandlers(mux *http.ServeMux) {
	mux.Handle("/get", TraceMiddleware(http.HandlerFunc(ListTodoHandler)))
	mux.Handle("/create", TraceMiddleware(http.HandlerFunc(AddTodoHandler)))
	mux.Handle("/update", TraceMiddleware(http.HandlerFunc(UpdateTodoHandler)))
	mux.Handle("/delete", TraceMiddleware(http.HandlerFunc(DeleteTodoHandler)))

	// Serve static /about page
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve dynamic list page
	mux.Handle("/list", TraceMiddleware(http.HandlerFunc(ListHandler)))
}
