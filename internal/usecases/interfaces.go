package usecases

import (
	"context"
	"prices/internal/domain"
)

type (
	ItemsRepository interface {
		GetActualPrice(ctx context.Context, itemName string) (float64, string, error)
	}

	ItemsUsecases interface {
		PriceIsHigherThanActual(ctx context.Context, item *domain.Item) (bool, error)
	}
)
