package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type PaymentRepository interface {
	ListActivePlans() ([]domain.Plan, error)
	FindPlanByID(id int64) (*domain.Plan, error)
	FindPlanByName(name string) (*domain.Plan, error)
	FindLatestCompletedPaymentByUserID(userID string) (*domain.Payment, error)
	CreatePayment(payment domain.Payment) (*domain.Payment, error)
	UpdatePaymentStatusByReference(referenceID, status string) error
}
