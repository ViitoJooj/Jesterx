package repository

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(orderID string) (*domain.Order, error)
	ListBySite(websiteID string) ([]domain.Order, error)
	ListSince(from, to time.Time) ([]domain.Order, error)
	UpdateStatus(orderID string, status domain.OrderStatus) error
}
