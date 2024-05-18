package main

import (
	"context"
	"log"
	"time"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/app/repository/orders_storage"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/app/server"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/app/services/warehouses_management_system"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/app/usecases/orders_management_system"
	middleware "github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/internal/middleware/errors"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/pkg/postgres"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/pkg/transaction_manager"
	"google.golang.org/grpc"
	// "github.com/Balun-courses/microservices_like_in_bigtech/lecture_8/orders_management_system/pkg/transaction_manager"
)

const DSN = "user=user password=password host=localhost port=6532 dbname=orders_management_system sslmode=require pool_max_conns=10"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// repository
	pool, err := postgres.NewConnectionPool(ctx, DSN,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		log.Fatal(err)
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
		GRPCPort:        ":8082",
		GRPCGatewayPort: ":8080",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.New(ctx, config, server.Deps{ //Dependency injection
		OMSUsecase: omsUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
