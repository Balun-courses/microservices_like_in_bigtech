package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/usecases/orders_management_system"
	pb "github.com/moguchev/microservices_courcse/orders_management_system/pkg/api/orders_management_system"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/closer"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config - server config
type Config struct {
	GRPCPort        string
	GRPCGatewayPort string
	DebugPort       string

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

	healthchecks   []func() error
	healthchecksMx sync.Mutex

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	grpcGateway struct {
		lis    net.Listener
		server *http.Server
	}

	internal struct {
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

	// grpc 82
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

	// grpc gateway 80, 443
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/
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

	// internal 84
	{
		mux := http.NewServeMux()
		// k8s readiness probe
		mux.HandleFunc("/healthz", srv.healthcheckHandler)
		// prometheus metrics
		mux.Handle("/metrics", promhttp.Handler())
		// pprof
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// TODO: swagger UI

		httpServer := &http.Server{Handler: mux}
		lis, err := net.Listen("tcp", cfg.DebugPort)
		if err != nil {
			return nil, fmt.Errorf("server: failed to listen: %v", err)
		}

		srv.internal.server = httpServer
		srv.internal.lis = lis
	}

	return srv, nil
}

func (s *Server) AddHealthcheck(hc func() error) {
	s.healthchecksMx.Lock()
	defer s.healthchecksMx.Unlock()

	s.healthchecks = append(s.healthchecks, hc)
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	s.healthchecksMx.Lock()
	defer s.healthchecksMx.Unlock()

	for _, hc := range s.healthchecks {
		if err := hc(); err != nil {
			_ = encode(w, http.StatusInternalServerError, struct{ Error string }{Error: err.Error()})
			return
		}
	}

	_ = encode(w, http.StatusOK, struct{}{})
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// Run - serve
func (s *Server) Run(ctx context.Context) error {
	group := errgroup.Group{}

	group.Go(func() error {
		closer.Add(func(ctx context.Context) error {
			s.grpc.server.Stop()
			return nil
		})

		logger.Info(ctx, "start serve", s.grpc.lis.Addr())
		if err := s.grpc.server.Serve(s.grpc.lis); err != nil {
			return fmt.Errorf("server: serve grpc: %v", err)
		}
		return nil
	})

	group.Go(func() error {
		closer.Add(s.grpcGateway.server.Shutdown)

		logger.Info(ctx, "start serve", s.grpcGateway.lis.Addr())
		if err := s.grpcGateway.server.Serve(s.grpcGateway.lis); err != nil {
			return fmt.Errorf("server: serve grpc gateway: %v", err)
		}
		return nil
	})

	group.Go(func() error {
		closer.Add(s.internal.server.Shutdown)

		logger.Info(ctx, "start serve", s.internal.lis.Addr())
		if err := s.internal.server.Serve(s.internal.lis); err != nil {
			return fmt.Errorf("server: serve internal: %v", err)
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
