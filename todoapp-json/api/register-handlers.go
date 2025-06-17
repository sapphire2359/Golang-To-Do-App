package api

import "net/http"

// registering all handlers wrapped with middleware trace function
func RegisterHandlers(mux *http.ServeMux) {
	mux.Handle("/get", TraceMiddleware(http.HandlerFunc(getHandler)))
	mux.Handle("/create", TraceMiddleware(http.HandlerFunc(createHandler)))
	mux.Handle("/update", TraceMiddleware(http.HandlerFunc(updateHandler)))
	mux.Handle("/delete", TraceMiddleware(http.HandlerFunc(deleteHandler)))

	// Serve static /about page
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve dynamic list page
	mux.Handle("/list", TraceMiddleware(http.HandlerFunc(listHandler)))
}
