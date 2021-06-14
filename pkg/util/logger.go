package util

import (
	"context"

	log "sunflower/pkg/api/gen/log"

	"go.uber.org/zap"
	"goa.design/goa/v3/middleware"
)

func L(ctx context.Context, logger *log.Logger) *zap.SugaredLogger {
	reqID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if ok {
		return logger.With(zap.String("reqID", reqID))
	}

	return logger.SugaredLogger
}
