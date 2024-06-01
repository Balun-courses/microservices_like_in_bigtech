package errors

import (
	"context"
	stderrors "errors"

	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorsUnaryInterceptor - convert any arror to rpc error
func ErrorsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		if _, ok := status.FromError(err); ok {
			return
		}

		switch {
		case stderrors.Is(err, models.ErrAlreadyExists):
			err = status.Error(codes.AlreadyExists, err.Error())
		case stderrors.Is(err, models.ErrUnimplemented):
			err = status.Error(codes.Unimplemented, err.Error())
		default:
			err = status.Error(codes.Internal, err.Error())
		}

		return
	}
}
