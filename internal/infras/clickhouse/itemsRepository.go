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
	query := fmt.Sprintf(`select i.price, i.unit from items i where ngramDistance(name, ?) < ? as diff order by diff desc, created_at desc limit 1;`)
	var res struct {
		Price float64 `ch:"price"`
		Unit  string  `ch:"unit"`
	}
	if err := r.client.QueryRow(ctx, query, itemName, diff).ScanStruct(&res); err != nil {
		return 0, "", err
	}

	return res.Price, res.Unit, nil
}