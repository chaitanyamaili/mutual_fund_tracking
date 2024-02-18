package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type API struct {
	mux   *httptreemux.ContextMux
	otmux http.Handler
	mw    []Middleware
}

func NewAPI(mw ...Middleware) *API {

	// Create an OpenTelemetry HTTP Handler which wraps our router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if a client request includes the appropriate headers
	// https://w3c.github.io/trace-context/

	mux := httptreemux.NewContextMux()

	return &API{
		mux:   mux,
		otmux: otelhttp.NewHandler(mux, "request"),
		mw:    mw,
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

func (a *API) Handle(method string, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler
	handler = wrapMiddleware(mw, handler)

	// Add the api's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// Execute each specific request
	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := ContextValues{
			// TracerUID: span.SpanContext().TraceID().String(),
			Now: time.Now(),
		}
		ctx = context.WithValue(ctx, key, &v)

		// Register this path and tracer uid for metrics later on
		_ = SetPath(ctx, path)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			// a.SignalShutdown()
			return
		}
	}

	a.mux.Handle(method, path, h)
}

func (a *API) NotFound(log *slog.Logger) {
	a.mux.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusNotFound
		url := r.URL.String()

		http.Error(w, fmt.Sprintf("path not found: %s", url), statusCode)
	}
}
