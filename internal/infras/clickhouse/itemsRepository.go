package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"prices/internal/usecases"
)

type ItemsRepository struct {
	client driver.Conn
}

func NewItemsRepository(client driver.Conn) usecases.ItemsRepository {
	return &ItemsRepository{client: client}
}

func (r *ItemsRepository) GetActualPrice(ctx context.Context, itemName string) (float64, string, error) {
	diff := 0.99
	query := fmt.Sprintf(`select i.price, i.unit from items i order by ngramSearchCaseInsensitive(name, ?) as diff desc, created_at desc limit 1`)
	var res struct {
		Price float64 `ch:"price"`
		Unit  string  `ch:"unit"`
	}
	if err := r.client.QueryRow(ctx, query, itemName, diff).ScanStruct(&res); err != nil {
		return 0, "", err
	}

	return res.Price, res.Unit, nil
}

func (r *ItemsRepository) MatchUnit(ctx context.Context, currentUnit string) (string, error) {
	query := fmt.Sprintf(`select unit from units order by ngramDistance(alt, ?) limit 1;`)
	var actualUnit string
	if err := r.client.QueryRow(ctx, query, currentUnit).Scan(&actualUnit); err != nil {
		return "", err
	}

	return actualUnit, nil
}
