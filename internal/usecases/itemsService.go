package usecases

import (
	"context"
	"fmt"
	"log"
	"math"
	"prices/internal/domain"
)

var (
	HighPriceError = fmt.Errorf("item overpriced")
)

type itemUsecases struct {
	repository ItemsRepository
}

func NewItemsUsecases(repository ItemsRepository) ItemsUsecases {
	return &itemUsecases{
		repository: repository,
	}
}

func (u *itemUsecases) PriceIsHigherThanActual(ctx context.Context, item *domain.Item) (uint, error) {
	actualPrice, actualUnit, err := u.repository.GetActualPrice(ctx, item.Product)
	if err != nil {
		return 0, err
	}

	currentUnit, err := u.repository.MatchUnit(ctx, item.Measure.Unit)
	if err != nil {
		return 0, err
	}
	log.Println("matched unit is ", currentUnit)
	log.Println("current price is ", item.Price)
	log.Println("actual price is ", actualPrice)
	log.Println("actual unit is ", actualUnit)

	isHigher := false
	var overprice uint
	switch actualUnit {
	case "кг":
		isHigher, overprice = u.proccessKg(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	case "десяток":
		isHigher, overprice = u.processDec(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	case "пачка":
		isHigher, overprice = u.processPack(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	case "л":
		isHigher, overprice = u.processLiter(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	case "рул":
		isHigher, overprice = u.processRul(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	case "шт.":
		isHigher, overprice = u.processOne(actualPrice, currentUnit, item.Price, item.Measure.Amount)
		break
	}

	if isHigher {
		return overprice, HighPriceError
	}

	return 0, nil
}

func (u *itemUsecases) proccessKg(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	var currentPrice float64

	switch currentUnit {
	case "г":
		currentPrice = price * (1000 / amount)
		break
	case "кг":
		currentPrice = price * (1 / amount)
		break
	default:
		currentPrice = price / amount
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0
}

func (u *itemUsecases) processDec(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	var currentPrice float64

	switch currentUnit {
	case "шт":
		currentPrice = price * 10
		break
	default:
		currentPrice = price
		break
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0
}

func (u *itemUsecases) processPack(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	var currentPrice float64

	switch currentUnit {
	default:
		currentPrice = price / amount
		break
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0
}

func (u *itemUsecases) processLiter(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	var currentPrice float64

	switch currentUnit {
	default:
		currentPrice = price * amount
		break
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0
}

func (u *itemUsecases) processRul(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	var currentPrice float64

	switch currentUnit {
	default:
		currentPrice = price / amount
		break
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0
}

func (u *itemUsecases) processOne(actualPrice float64, currentUnit string, price float64, amount float64) (bool, uint) {
	log.Println("proccessing one")
	var currentPrice float64

	switch currentUnit {
	default:
		currentPrice = price
		break
	}

	if currentPrice > actualPrice {
		return true, u.calcOverprice(currentPrice, actualPrice)
	}

	return false, 0

}

func (u *itemUsecases) calcOverprice(currentPrice, actualPrice float64) uint {
	overprice := uint(math.Floor((currentPrice-actualPrice)/actualPrice) * 100)
	log.Println("current ", currentPrice, " actual ", actualPrice, " overprice ", overprice, "%")
	return overprice
}
