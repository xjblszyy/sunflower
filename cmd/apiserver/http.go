package apiserver

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"

	"sunflower/pkg/middleware/metrics"

	uploadscoresvr "sunflower/pkg/api/gen/http/score/server"

	"sunflower/pkg/api/gen/log"
	"sunflower/pkg/api/gen/score"

	mdlwr "sunflower/pkg/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, addr string,
	uploadScoreEndpoints *score.Endpoints,
	wg *sync.WaitGroup, errc chan error, logger *log.Logger, metrics *metrics.Prometheus, debug bool) {
	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = logger
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		uploadServer *uploadscoresvr.Server
	)
	{
		eh := errorHandler(logger)
		uploadServer = uploadscoresvr.New(uploadScoreEndpoints, mux, dec, enc, eh, mdlwr.GoaErrorFormatterFunc)
	}
	// Configure the mux.
	uploadscoresvr.Mount(mux, uploadServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		if debug {
			handler = httpmdlwr.Debug(mux, os.Stdout)(handler)
		}

		handler = mdlwr.PopulateRequestContext()(handler)
		handler = httpmdlwr.RequestID()(handler)

		if metrics != nil {
			handler = metrics.HandlerFunc(adapter)(handler)
		} else {
			handler = httpmdlwr.Log(adapter)(handler)
		}
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: addr, Handler: handler}

	for _, m := range uploadServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)

	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Infof("HTTP server listening on %q", addr)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Infof("shutting down HTTP server at %q", addr)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.With(zap.String("id", id)).Error(err.Error())
	}
}
