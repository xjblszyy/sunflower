package apiserver

import (
	"context"
	"net"
	"sync"

	"sunflower/pkg/api/gen/log"
	"sunflower/pkg/api/gen/score"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcmdlwr "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// handleGRPCServer starts configures and starts a gRPC server on the given
// URL. It shuts down the server if any error is received in the error channel.
func HandleGRPCServer(ctx context.Context, addr string,
	uploadScoreEndpoints *score.Endpoints,
	wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {
	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to gRPC requests and
	// responses.

	// adapter := middleware.NewLogger(logger)

	// Initialize gRPC server with the middleware.
	srv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcmdlwr.UnaryRequestID(),
			grpcmdlwr.UnaryServerLog(logger),
		),
	)

	// Register the servers.

	for svc, info := range srv.GetServiceInfo() {
		for _, m := range info.Methods {
			logger.SugaredLogger.Infof("serving gRPC method %s", svc+"/"+m.Name)
		}
	}

	// Register the server reflection service on the server.
	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
	reflection.Register(srv)

	(*wg).Add(1)

	go func() {
		defer (*wg).Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				errc <- err
			}

			logger.SugaredLogger.Infof("gRPC server listening on %q", addr)
			errc <- srv.Serve(lis)
		}()

		<-ctx.Done()
		logger.SugaredLogger.Infof("shutting down gRPC server at %q", addr)
		srv.Stop()
	}()
}
