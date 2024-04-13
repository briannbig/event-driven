package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/briannbig/event-driven/internal/application/controller"
	"github.com/briannbig/event-driven/internal/application/usecase"
	"github.com/briannbig/event-driven/internal/domain/event"
	"github.com/briannbig/event-driven/internal/infra/queue"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx := context.Background()

	// initialize queue_
	queue_ := queue.NewQueueAdapter(nats.DefaultURL)

	

	// use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(queue_)
	processPaymentUseCase := usecase.NewProcessOrderPaymentUseCase(queue_)
	stockMovementUseCase := usecase.NewStockMovementUseCase()
	sendOrderEmailUseCase := usecase.NewSendOrderEmailUseCase()

	// controllers
	orderController := controller.NewOrderController(createOrderUseCase, processPaymentUseCase, stockMovementUseCase, sendOrderEmailUseCase)

	// register routes
	http.HandleFunc("POST /create-order", orderController.CreateOrder_)

	// mapping listeners
	var list map[reflect.Type][]func(payload *nats.Msg) = map[reflect.Type][]func(payload *nats.Msg){
		reflect.TypeOf(event.OrderCreatedEvent{}): {
			orderController.ProcessOrderPayment,
			orderController.StockMovement,
			orderController.SendOrderEmail,
		},
	}

	// register listeners
	for eventType, handlers := range list {
		for _, handler := range handlers {
			queue_.ListenerRegister(eventType, handler)
		}
	}

	// connect queue
	err := queue_.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connect queue %s", err)
	}
	defer queue_.Disconnect(ctx)

	// start consuming queues
	OrderCreatedEvent := reflect.TypeOf(event.OrderCreatedEvent{}).Name()

	go func(ctx context.Context, queueName string) {
		err = queue_.StartConsuming(ctx, queueName)
		if err != nil {
			log.Fatalf("Error running consumer %s: %s", queueName, err)
		}
	}(ctx, OrderCreatedEvent)

	// start server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
