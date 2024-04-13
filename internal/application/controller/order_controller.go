package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/briannbig/event-driven/internal/application/dto"
	"github.com/briannbig/event-driven/internal/application/usecase"
	"github.com/briannbig/event-driven/internal/domain/event"
	"github.com/nats-io/nats.go"
)

type OrderController struct {
	createOrderUserCase        *usecase.CreateOrderUseCase
	processOrderPaymentUseCase *usecase.ProcessOrderPaymentUseCase
	stockMovementUseCase       *usecase.StockMovementUseCase
	sendOrderEmailUseCase      *usecase.SendOrderEmailUseCase
}

func NewOrderController(createOrderUserCase *usecase.CreateOrderUseCase,
	processOrderPaymentUseCase *usecase.ProcessOrderPaymentUseCase,
	stockMovementUseCase *usecase.StockMovementUseCase,
	sendOrderEmailUseCase *usecase.SendOrderEmailUseCase) *OrderController {
	return &OrderController{
		createOrderUserCase:        createOrderUserCase,
		processOrderPaymentUseCase: processOrderPaymentUseCase,
		stockMovementUseCase:       stockMovementUseCase,
		sendOrderEmailUseCase:      sendOrderEmailUseCase,
	}
}

func (u *OrderController) CreateOrder(payload *nats.Msg) {
	var requestData dto.CreateOrderDTO
	err := json.Unmarshal(payload.Data, &requestData)
	if err != nil {
		log.Printf("error unmarshalling payload --- %s", err.Error())
		return
	}
	err = u.createOrderUserCase.Execute(requestData)
	if err != nil {
		return
	}

}

func (u *OrderController) CreateOrder_(w http.ResponseWriter, r *http.Request) {
	var requestData dto.CreateOrderDTO
	json.NewDecoder(r.Body).Decode(&requestData)
	err := u.createOrderUserCase.Execute(requestData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (u *OrderController) ProcessOrderPayment(payload *nats.Msg) {
	var body event.OrderCreatedEvent
	err := json.Unmarshal(payload.Data, &body)
	if err != nil {
		log.Printf("error unmarshalling payload --- %s", err.Error())
		return
	}
	err = u.processOrderPaymentUseCase.Execute(&body)
	if err != nil {
		return
	}
}

func (u *OrderController) StockMovement(payload *nats.Msg) {
	var body event.OrderCreatedEvent
	err := json.Unmarshal(payload.Data, &body)
	if err != nil {
		log.Printf("error unmarshalling payload --- %s", err.Error())
		return
	}
	err = u.stockMovementUseCase.Execute(&body)
	if err != nil {
		return
	}
}

func (u *OrderController) SendOrderEmail(payload *nats.Msg) {
	var body event.OrderCreatedEvent
	err := json.Unmarshal(payload.Data, &body)
	if err != nil {
		log.Printf("error unmarshalling payload --- %s", err.Error())
		return
	}
	err = u.sendOrderEmailUseCase.Execute(&body)
	if err != nil {
		return
	}
}
