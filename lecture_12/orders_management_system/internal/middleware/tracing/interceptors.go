package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	traceIDKey = "x-trace-id"
)

// DebugOpenTracingUnaryServerInterceptor - ...
func DebugOpenTracingUnaryServerInterceptor(logRequest, logResponse bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			span, ctx = opentracing.StartSpanFromContext(ctx, info.FullMethod)
			defer span.Finish()

			spanContext, ok := span.Context().(jaeger.SpanContext)
			if ok {
				ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIDKey, spanContext.TraceID().String()))

				header := metadata.New(map[string]string{traceIDKey: spanContext.TraceID().String()})
				err := grpc.SendHeader(ctx, header)
				if err != nil {
					return nil, err
				}
			}
		}

		if pbMsg, ok := req.(proto.Message); ok && logRequest {
			if jsonRequest, err := protojson.Marshal(pbMsg); err == nil {
				span.LogKV("grpc_request", string(jsonRequest))
			}
		}

		res, err := handler(ctx, req)
		if err != nil {
			ext.Error.Set(span, true) // тоже самое что span.SetTag("error", true)
			span.LogKV("grpc_error", err)
		} else {
			if pbMsg, ok := res.(proto.Message); ok && logResponse {
				if jsonResponse, err := protojson.Marshal(pbMsg); err == nil {
					span.LogKV("grpc_response", string(jsonResponse))
				}
			}
		}

		return res, err
	}
}
