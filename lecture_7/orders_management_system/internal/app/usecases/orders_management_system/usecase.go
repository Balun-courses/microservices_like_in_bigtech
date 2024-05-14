package orders_management_system

import (
	"context"
	"errors"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_7/orders_management_system/internal/app/models"
)

var (
	// ErrReserveStocks - ...
	ErrReserveStocks = errors.New("failed to reserve stock")
)

// UsecaseInterface - интерфейс бизнес логики
type UsecaseInterface interface {
	// CreateOrder - создание заказа
	//
	// @errors: ErrReserveStocks
	CreateOrder(ctx context.Context, userID models.UserID, info CreateOrderInfo) (*models.Order, error)
}

// Бизнес логика не зависит ни от чего кроме доменных моделей!
// Объявляем интерфейсы зависимостей в месте использования!
// Задаем контракт поведения для адаптеров (порты)

//go:generate mockery --name=WarehouseManagementSystem --filename=warehouse_management_system_mock.go --disable-version-string
//go:generate mockery --name=OrdersStorage --filename=orders_storage_mock.go --disable-version-string

type (
	// WarehouseManagementSystem - то что отвечает за резервирование товаров на складе
	WarehouseManagementSystem interface {
		// ReserveStocks - резервация стоков на складах
		ReserveStocks(ctx context.Context, userID models.UserID, items []models.Item) error
	}

	// OrdersStorage - репозиторий сервиса OMS
	OrdersStorage interface {
		// CreateOrder - создание записи заказа в БД
		//
		// @errors: models.ErrAlreadyExists
		CreateOrder(ctx context.Context, order *models.Order) error
	}
)

// Deps - зависимости нашего usecase
type Deps struct {
	WarehouseManagementSystem
	OrdersStorage
}

// usecase - реализация
type usecase struct {
	Deps
}

// NewUsecase - возвращаем реализацию UsecaseInterface
func NewUsecase(d Deps) UsecaseInterface {
	return &usecase{
		Deps: d,
	}
}
