//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/wire"
	"log"
	"prices/internal/broker"
	"prices/internal/config"
	"prices/internal/events"
	ch "prices/internal/infras/clickhouse"
	"prices/internal/usecases"
	"time"
)

func Init() (*App, func(), error) {
	panic(
		wire.Build(NewApplication,
			wire.NewSet(config.New),
			wire.NewSet(initDB),
			wire.NewSet(initBroker),
			wire.NewSet(ch.NewItemsRepository),
			wire.NewSet(usecases.NewItemsUsecases),

			wire.NewSet(events.NewValidationHandler)))
}

func initDB(cfg *config.Config) (driver.Conn, func(), error) {
	log.Println(cfg.DB)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.DB.Host, cfg.DB.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.DB.Name,
			Username: cfg.DB.User,
			Password: cfg.DB.Pass,
		},
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, func() {
			conn.Close()
		}, fmt.Errorf("cannot ping database: ", err)
	}

	return conn, func() {
		conn.Close()
	}, nil
}

func initBroker(cfg *config.Config) (broker.MessageBroker, func(), error) {
	mb, err := broker.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	return mb, func() {
		mb.Close()
	}, nil
}
