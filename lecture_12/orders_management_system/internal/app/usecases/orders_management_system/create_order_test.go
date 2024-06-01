//go:build test

package orders_management_system

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/models"
	"github.com/moguchev/microservices_courcse/orders_management_system/internal/app/usecases/orders_management_system/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_usecase_CreateOrder(t *testing.T) {
	var (
		ctx  = context.Background() // dummy
		date = time.Now()
	)
	type fields struct {
		WarehouseManagementSystem *mocks.WarehouseManagementSystem
		OrdersStorage             *mocks.OrdersStorage
	}

	type args struct {
		ctx    context.Context
		userID models.UserID
		info   CreateOrderInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Order
		wantErr bool

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx:    ctx, // dumm
				userID: 1,
				info: CreateOrderInfo{
					Items: []models.Item{
						{
							SKU:         models.SKU{ID: 2, Name: "Item 2"},
							Quantity:    3,
							WarehouseID: 4,
						},
					},
					DeliveryOrderInfo: models.DeliveryOrderInfo{
						DeliveryVariantID: 5,
						DeliveryDate:      date,
					},
				},
			},
			want: &models.Order{
				UserID: 1,
				Items: []models.Item{
					{
						SKU:         models.SKU{ID: 2, Name: "Item 2"},
						Quantity:    3,
						WarehouseID: 4,
					},
				},
				DeliveryOrderInfo: models.DeliveryOrderInfo{
					DeliveryVariantID: 5,
					DeliveryDate:      date,
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.WarehouseManagementSystem.On("ReserveStocks", ctx, models.UserID(1), []models.Item{
					{
						SKU:         models.SKU{ID: 2, Name: "Item 2"},
						Quantity:    3,
						WarehouseID: 4,
					},
				}).
					Return(nil)
				f.OrdersStorage.On("CreateOrder", ctx, mock.MatchedBy(func(order *models.Order) bool {
					return order != nil &&
						order.UserID == 1 &&
						reflect.DeepEqual(order.Items, []models.Item{
							{
								SKU:         models.SKU{ID: 2, Name: "Item 2"},
								Quantity:    3,
								WarehouseID: 4,
							},
						}) &&
						reflect.DeepEqual(
							order.DeliveryOrderInfo, models.DeliveryOrderInfo{
								DeliveryVariantID: 5,
								DeliveryDate:      date,
							},
						) &&
						order.ID != models.OrderID{} // not empty
				})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.WarehouseManagementSystem.AssertNumberOfCalls(t, "ReserveStocks", 1)
				f.OrdersStorage.AssertNumberOfCalls(t, "CreateOrder", 1)
			},
		},
		{
			name: "Test 2. Negative. WarehouseManagementSystem returns error.",
			args: args{
				ctx:    ctx, // dumm
				userID: 1,
				info: CreateOrderInfo{
					Items: []models.Item{
						{
							SKU:         models.SKU{ID: 2, Name: "Item 2"},
							Quantity:    3,
							WarehouseID: 4,
						},
					},
					DeliveryOrderInfo: models.DeliveryOrderInfo{
						DeliveryVariantID: 5,
						DeliveryDate:      date,
					},
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.WarehouseManagementSystem.On("ReserveStocks", ctx, models.UserID(1), []models.Item{
					{
						SKU:         models.SKU{ID: 2, Name: "Item 2"},
						Quantity:    3,
						WarehouseID: 4,
					},
				}).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.WarehouseManagementSystem.AssertNumberOfCalls(t, "ReserveStocks", 1)
				f.OrdersStorage.AssertNumberOfCalls(t, "CreateOrder", 0)
			},
		},
		{
			name: "Test 3. Negative. CreateOrder returns error.",
			args: args{
				ctx:    ctx, // dumm
				userID: 1,
				info: CreateOrderInfo{
					Items: []models.Item{
						{
							SKU:         models.SKU{ID: 2, Name: "Item 2"},
							Quantity:    3,
							WarehouseID: 4,
						},
					},
					DeliveryOrderInfo: models.DeliveryOrderInfo{
						DeliveryVariantID: 5,
						DeliveryDate:      date,
					},
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.WarehouseManagementSystem.On("ReserveStocks", ctx, models.UserID(1), []models.Item{
					{
						SKU:         models.SKU{ID: 2, Name: "Item 2"},
						Quantity:    3,
						WarehouseID: 4,
					},
				}).
					Return(nil)
				f.OrdersStorage.On("CreateOrder", ctx, mock.MatchedBy(func(order *models.Order) bool {
					return order != nil &&
						order.UserID == 1 &&
						reflect.DeepEqual(order.Items, []models.Item{
							{
								SKU:         models.SKU{ID: 2, Name: "Item 2"},
								Quantity:    3,
								WarehouseID: 4,
							},
						}) &&
						reflect.DeepEqual(
							order.DeliveryOrderInfo, models.DeliveryOrderInfo{
								DeliveryVariantID: 5,
								DeliveryDate:      date,
							},
						) &&
						order.ID != models.OrderID{} // not empty
				})).
					Return(models.ErrAlreadyExists)
			},
			assert: func(t *testing.T, f *fields) {
				f.WarehouseManagementSystem.AssertNumberOfCalls(t, "ReserveStocks", 1)
				f.OrdersStorage.AssertNumberOfCalls(t, "CreateOrder", 3)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				WarehouseManagementSystem: mocks.NewWarehouseManagementSystem(t),
				OrdersStorage:             mocks.NewOrdersStorage(t),
			}
			oms := &usecase{
				Deps: Deps{
					WarehouseManagementSystem: f.WarehouseManagementSystem,
					OrdersStorage:             f.OrdersStorage,
				},
			}
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := oms.CreateOrder(tt.args.ctx, tt.args.userID, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// assert
			if got != nil { // зануляем так как не можем проверить
				got.ID = models.OrderID{}
			}
			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
