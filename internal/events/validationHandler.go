package events

import (
	"context"
	"encoding/json"
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

	ok, err := h.uc.PriceIsHigherThanActual(context.Background(), &item)
	if err != nil {
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
	if !ok {
		data, _ := json.Marshal(&errorReq{ErrorMessage: "Цена не превышает максимально допустимую"})
		h.broker.Publish(&nats.Msg{
			Subject: h.cfg.Nats.Queues.Errors,
			Header:  msg.Header,
			Data:    data,
		})
		return
	}
	req := &statusReq{OperationType: "validation"}
	data, _ := json.Marshal(&req)
	if err := h.broker.Publish(&nats.Msg{
		Subject: h.cfg.Nats.Queues.Status,
		Header:  msg.Header,
		Data:    data,
	}); err != nil {
		log.Println("error while publishing")
		return
	}
}
