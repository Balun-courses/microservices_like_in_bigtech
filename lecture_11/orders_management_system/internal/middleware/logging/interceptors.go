package logging

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_11/orders_management_system/pkg/logger"
	"google.golang.org/grpc"
)

// LogErrorUnaryInterceptor - log interceptor
func LogErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logCtx := logger.ToContext(context.Background(),
			logger.FromContext(ctx).With(
				"operation", info.FullMethod,
				"component", "middleware",
			),
		)

		logger.Debug(logCtx, "receive request")
		resp, err = handler(ctx, req)
		logger.Debug(logCtx, "handle request")

		if err != nil {
			// 4ХХ -> warn
			// 5ХХ -> Error
			logger.Error(logCtx, err.Error())
		}

		return
	}
}
