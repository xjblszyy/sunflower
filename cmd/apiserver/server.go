package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"sunflower/pkg/libs/panichandler"

	"sunflower/pkg/api/gen/score"

	"sunflower/config"
	"sunflower/pkg/api/gen/log"
	"sunflower/pkg/app/apiserver/handler"

	"go.uber.org/zap"

	metricsMlwr "sunflower/pkg/middleware/metrics"
)

func RunServer(cfg *config.Config, metrics *metricsMlwr.Prometheus) {
	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(cfg.ServiceName, !cfg.Debug)
	}

	// Initialize the services.
	var (
		scoreSvc score.Service
	)
	{
		scoreSvc = handler.NewScore(logger)

	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		uploadEndpoints *score.Endpoints
	)
	{
		uploadEndpoints = score.NewEndpoints(scoreSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	// serve http
	handleHTTPServer(ctx,
		cfg.Server.HttpAddr,
		uploadEndpoints,
		&wg, errc, logger, metrics, cfg.Debug)

	// serve grpc
	HandleGRPCServer(ctx,
		cfg.Server.GrpcAddr,
		uploadEndpoints,
		&wg, errc, logger, cfg.Debug)

	// Wait for signal.
	logger.Infof("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Info("exited")
}

// 开启 pprof
func RunDebugPprofServer(addr string) {
	defer panichandler.ZapHandler(zap.L()).Handle()
	zap.L().Sugar().Infof("启动 pprof 监听 %s.", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		zap.L().Error("开启 pprof 监听失败 %s", zap.Error(err))
	}
}
