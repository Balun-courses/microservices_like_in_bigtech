package main

import (
	"context"
	"log"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/repository/orders_storage"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/server"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/services/warehouses_management_system"
	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/usecases/orders_management_system"
	middleware "github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/middleware/errors"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// repository

	storage := orders_storage.New()

	// services

	wmsClient := warehouses_management_system.NewClient()

	// usecases

	omsUsecase := orders_management_system.NewUsecase(orders_management_system.Deps{ // Dependency injection
		WarehouseManagementSystem: wmsClient,
		ReadOrdersStorage:         storage,
		WriteOrdersStorage:        storage,
	})

	// delivery
	config := server.Config{
		GRPCPort:        ":8082",
		GRPCGatewayPort: ":8080",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.New(ctx, config, server.Deps{ // Dependency injection
		OMSUsecase: omsUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
