package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/usecases/orders_management_system"
	pb "github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/pkg/api/orders_management_system"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config - server config
type Config struct {
	GRPCPort        string
	GRPCGatewayPort string

	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	OMSUsecase orders_management_system.UsecaseInterface
}

// Server
type Server struct {
	pb.UnimplementedOrdersManagementSystemServiceServer
	Deps

	validator *protovalidate.Validator

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	grpcGateway struct {
		lis    net.Listener
		server *http.Server
	}
}

// New - returns *Server
func New(ctx context.Context, cfg Config, d Deps) (*Server, error) {
	srv := &Server{
		Deps: d,
	}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				// Добавляем сюда все запросы наши
				&pb.CreateOrderRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	// grpc
	{
		grpcServerOptions := unaryInterceptorsToGrpcServerOptions(cfg.UnaryInterceptors...)
		grpcServerOptions = append(grpcServerOptions,
			grpc.ChainUnaryInterceptor(cfg.ChainUnaryInterceptors...),
		)

		grpcServer := grpc.NewServer(grpcServerOptions...)
		pb.RegisterOrdersManagementSystemServiceServer(grpcServer, srv)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			return nil, fmt.Errorf("server: failed to listen: %v", err)
		}

		srv.grpc.server = grpcServer
		srv.grpc.lis = lis
	}

	// grpc gateway
	{
		mux := runtime.NewServeMux()
		if err := pb.RegisterOrdersManagementSystemServiceHandlerServer(ctx, mux, srv); err != nil {
			return nil, fmt.Errorf("server: failed to register handler: %v", err)
		}

		httpServer := &http.Server{Handler: mux}

		lis, err := net.Listen("tcp", cfg.GRPCGatewayPort)
		if err != nil {
			return nil, fmt.Errorf("server: failed to listen: %v", err)
		}

		srv.grpcGateway.server = httpServer
		srv.grpcGateway.lis = lis
	}

	return srv, nil
}

// Run - serve
func (s *Server) Run(ctx context.Context) error {
	group := errgroup.Group{}

	group.Go(func() error {
		log.Println("start serve", s.grpc.lis.Addr())
		if err := s.grpc.server.Serve(s.grpc.lis); err != nil {
			return fmt.Errorf("server: serve grpc: %v", err)
		}
		return nil
	})

	group.Go(func() error {
		log.Println("start serve", s.grpcGateway.lis.Addr())
		if err := s.grpcGateway.server.Serve(s.grpcGateway.lis); err != nil {
			return fmt.Errorf("server: serve grpc gateway: %v", err)
		}
		return nil
	})

	return group.Wait()
}

func unaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
