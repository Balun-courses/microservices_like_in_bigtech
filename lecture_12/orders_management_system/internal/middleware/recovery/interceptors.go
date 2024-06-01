package recovery

import (
	"context"
	"runtime/debug"

	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoverUnaryInterceptor - recover panics
func RecoverUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {
		defer func() {
			if v := recover(); v != nil {
				logger.ErrorKV(ctx, "recover panic",
					"panic", v,
					"stacktrace", string(debug.Stack()),
					"operation", info.FullMethod,
					"component", "middleware",
				)

				err = status.Error(codes.Internal, codes.Internal.String()) // return error
			}
		}()

		return handler(ctx, req)
	}
}
