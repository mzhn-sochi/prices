package usecases

import (
	"context"
	"fmt"
	"prices/internal/domain"
)

var (
	HighPriceError = fmt.Errorf("item has high price")
)

type itemUsecases struct {
	repository ItemsRepository
}

func NewItemsUsecases(repository ItemsRepository) ItemsUsecases {
	return &itemUsecases{
		repository: repository,
	}
}

func (u *itemUsecases) PriceIsHigherThanActual(ctx context.Context, item *domain.Item) (bool, error) {
	diff := 5.0
	actualPrice, _, err := u.repository.GetActualPrice(ctx, item.Product)
	if err != nil {
		return false, err
	}

	if item.Price > (actualPrice + diff) {
		return true, nil
	}

	return false, nil
}
