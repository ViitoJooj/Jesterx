package domain

import (
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain/entities/enums"
	"github.com/google/uuid"
)

type JesterxPlans struct {
	UUID            uuid.UUID
	Name            string
	Description     string
	MaxWebsites     int
	MaxRouters      int
	MaxProducts     int
	CostPerSaleRate int
	Coin            enums.CoinType
	Price           int
	UpdatedAt       *time.Time
	CreatedAt       time.Time
}

func NewPlan(name string, description string, maxWebsites int, maxRouters int, maxProducts int, costPerSaleRate int, coin string, price int) (*JesterxPlans, error) {

	if name == "" {
		return nil, errors.New("Name cannot be null.")
	}

	coinType := enums.CoinType(coin)
	if coinType != enums.BRCoin && coinType != enums.EUACoin && coinType != enums.EURCoin {
		return nil, errors.New("Coin must be 'BRL', 'USD' or 'EUR'.")
	}

	if price <= 0 {
		return nil, errors.New("Price must be greater than 0.")
	}

	return &JesterxPlans{
		UUID:            uuid.Nil,
		Name:            name,
		Description:     description,
		MaxWebsites:     maxWebsites,
		MaxRouters:      maxRouters,
		MaxProducts:     maxProducts,
		CostPerSaleRate: costPerSaleRate,
		Coin:            coinType,
		Price:           price,
	}, nil
}
