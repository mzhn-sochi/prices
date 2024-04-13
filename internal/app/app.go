package app

import (
	"prices/internal/broker"
	"prices/internal/events"
)

type App struct {
	messageBroker     broker.MessageBroker
	validationHandler events.ValidationHandler
}

func NewApplication(broker broker.MessageBroker, validationHandler events.ValidationHandler) *App {
	return &App{
		messageBroker:     broker,
		validationHandler: validationHandler,
	}
}

func (a *App) Run() {
	if err := a.messageBroker.Consume(a.validationHandler.Handle); err != nil {
		panic(err)
	}
}
