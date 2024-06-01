package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MetricsUnaryInterceptor - ...
func MetricsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ interface{}, err error) {
		start := time.Now()
		defer func() {
			switch status.Code(err) {
			case codes.OK:
			case codes.Internal:
				// ++
				//
			}
			responseTimeHistogramObserve(info.FullMethod, err, time.Since(start))
		}()

		return handler(ctx, req)
	}
}
