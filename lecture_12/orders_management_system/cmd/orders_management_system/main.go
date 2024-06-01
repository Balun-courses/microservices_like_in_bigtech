package main

import (
	"context"
	"os"
	"time"

	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/repository/orders_storage"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/server"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/services/warehouses_management_system"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/usecases/orders_management_system"
	middleware_errors "github.com/moguchev/microservices_courcse/orders_management_system/internal/middleware/errors"
	middleware_logging "github.com/moguchev/microservices_courcse/orders_management_system/internal/middleware/logging"
	middleware_metrics "github.com/moguchev/microservices_courcse/orders_management_system/internal/middleware/metrics"
	middleware_recovery "github.com/moguchev/microservices_courcse/orders_management_system/internal/middleware/recovery"
	middleware_tracing "github.com/moguchev/microservices_courcse/orders_management_system/internal/middleware/tracing"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/logger"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/postgres"
	jaeger_tracing "github.com/moguchev/microservices_courcse/orders_management_system/pkg/tracing"
	"github.com/moguchev/microservices_courcse/orders_management_system/pkg/transaction_manager"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.SetLevel(zapcore.DebugLevel)

	logger.Info(ctx, "start app init")
	if err := jaeger_tracing.Init("orders-management-system"); err != nil {
		logger.Fatal(ctx, err)
	}

	// repository
	pool, err := postgres.NewConnectionPool(ctx, os.Getenv("DB_DSN"),
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		logger.Fatal(ctx, err)
	}

	txManager := transaction_manager.New(pool)

	storage := orders_storage.New(txManager)

	// services

	wmsClient := warehouses_management_system.NewClient()

	// usecases

	omsUsecase := orders_management_system.NewUsecase(orders_management_system.Deps{ // Dependency injection
		WarehouseManagementSystem: wmsClient,
		OrdersStorage:             storage,
		TransactionManager:        txManager,
	})

	// delivery
	config := server.Config{
		DebugPort:       os.Getenv("DEBUG_HTTP_PORT"), // ":8084"
		GRPCPort:        os.Getenv("GRPC_PORT"),       // ":8082"
		GRPCGatewayPort: os.Getenv("HTTP_PORT"),       // ":8080"
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			// https://github.com/grpc-ecosystem/go-grpc-middleware?tab=readme-ov-file#middleware
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer()), // Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
			middleware_logging.LogErrorUnaryInterceptor(),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true), // расширение для grpc_opentracing.OpenTracingServerInterceptor
			middleware_metrics.MetricsUnaryInterceptor(),
			middleware_recovery.RecoverUnaryInterceptor(), // можно использовать grpc_recovery
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware_errors.ErrorsUnaryInterceptor(), // далее наши остальные middleware
		},
	}

	srv, err := server.New(ctx, config, server.Deps{ // Dependency injection (DI)
		OMSUsecase: omsUsecase,
	})
	if err != nil {
		logger.Fatalf(ctx, "failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		logger.Errorf(ctx, "run: %v", err)
	}
}
