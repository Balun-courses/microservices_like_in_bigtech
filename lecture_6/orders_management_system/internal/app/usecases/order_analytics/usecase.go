package order_analytics

import (
	"context"

	"github.com/Balun-courses/microservices_like_in_bigtech/lecture_6/orders_management_system/internal/app/models"
)

// UsecaseInterface - интерфейс бизнес логики
type OrderAnalyticsUsecaseInterface interface {
	// SendOrderData - ...
	SendOrderData(ctx context.Context, order *models.Order) error
}

// Бизнес логика не зависит ни от чего кроме доменных моделей!
// Объявляем интерфейсы зависимостей в месте использования!
// Задаем контракт поведения для адаптеров (порты)
type (
	// WarehouseManagementSystem - то что отвечает за резервирование товаров на складе
	KafkaAdapter interface {
		// ReserveStocks - резервация стоков на складах
		ReserveStocks(ctx context.Context, userID models.UserID, items []models.Item) error
	}
)

// Deps - зависимости нашего usecase
type Deps struct {
}

// usecase - реализация
type usecase struct {
	Deps
}

// NewUsecase - возвращаем реализацию UsecaseInterface
func NewUsecase(d Deps) OrderAnalyticsUsecaseInterface {
	return &usecase{
		Deps: d,
	}
}

func (u *usecase) SendOrderData(ctx context.Context, order *models.Order) error {
	return models.ErrUnimplemented
}
