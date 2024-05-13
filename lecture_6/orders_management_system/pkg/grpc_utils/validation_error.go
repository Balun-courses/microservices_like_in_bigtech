package grpcutils

import (
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCValidationError - returns codes.InvalidArgument
func RPCValidationError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, stErr := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if stErr == nil {
			return st.Err()
		}
	}

	return status.Error(codes.InvalidArgument, err.Error())
}
