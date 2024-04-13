package events

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"prices/internal/broker"
	"prices/internal/config"
	"prices/internal/domain"
	"prices/internal/usecases"
	"strings"
)

type ValidationHandler interface {
	Handle(msg *nats.Msg)
}

type validationHandler struct {
	uc     usecases.ItemsUsecases
	broker broker.MessageBroker
	cfg    *config.Config
}

func NewValidationHandler(cfg *config.Config, usecases usecases.ItemsUsecases,
	broker broker.MessageBroker) ValidationHandler {
	return &validationHandler{cfg: cfg, uc: usecases, broker: broker}
}

type statusReq struct {
	OperationType string `json:"type"`
}

type overpriceReq struct {
	Percents uint `json:"percents"`
}

type errorReq struct {
	ErrorMessage string `json:"message"`
}

func (h *validationHandler) Handle(msg *nats.Msg) {
	var item domain.Item
	log.Printf("received new message %v\n", string(msg.Data))
	if err := json.Unmarshal(msg.Data, &item); err != nil {
		log.Println("error while unmarshal")
		return
	}

	ticketId := strings.TrimSpace(msg.Header.Get("ticket_id"))
	if ticketId == "" {
		log.Printf("empty ticket id\n")
		return
	}

	overprice, err := h.uc.PriceIsHigherThanActual(context.Background(), &item)
	if err != nil && !errors.Is(err, usecases.HighPriceError) {
		log.Println(item)
		log.Println("error while check price: ", err)
		data, _ := json.Marshal(&errorReq{ErrorMessage: "Товар отсутсвует в базе социальных"})
		h.broker.Publish(&nats.Msg{
			Subject: h.cfg.Nats.Queues.Errors,
			Header:  msg.Header,
			Data:    data,
		})
		return
	}
	if overprice == 0 {
		data, _ := json.Marshal(&errorReq{ErrorMessage: "Цена не превышает максимально допустимую"})
		h.broker.Publish(&nats.Msg{
			Subject: h.cfg.Nats.Queues.Errors,
			Header:  msg.Header,
			Data:    data,
		})
		return
	} else {
		oReq := &overpriceReq{Percents: overprice}
		oData, _ := json.Marshal(&oReq)
		if err := h.broker.Publish(&nats.Msg{
			Subject: h.cfg.Nats.Queues.Overprice,
			Header:  msg.Header,
			Data:    oData,
		}); err != nil {
			log.Println("error while publishing to overprice subject")
			return
		}

		sReq := &statusReq{OperationType: "validation"}
		sData, _ := json.Marshal(&sReq)
		if err := h.broker.Publish(&nats.Msg{
			Subject: h.cfg.Nats.Queues.Status,
			Header:  msg.Header,
			Data:    sData,
		}); err != nil {
			log.Println("error while publishing to validation subject")
			return
		}
	}
}
