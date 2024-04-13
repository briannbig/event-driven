package usecase

import (
	"context"
	"fmt"

	"github.com/briannbig/event-driven/internal/application/dto"
	"github.com/briannbig/event-driven/internal/domain/entity"
	"github.com/briannbig/event-driven/internal/domain/event"
	"github.com/briannbig/event-driven/internal/domain/queue"
)

type CreateOrderUseCase struct {
	publisher queue.Publisher
}

func NewCreateOrderUseCase(publisher queue.Publisher) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		publisher: publisher,
	}
}

func (u *CreateOrderUseCase) Execute(ctx context.Context, input dto.CreateOrderDTO) error {
	fmt.Println("--- CreateOrderUseCase ---")

	// create order
	order, err := entity.NewOrderEntity()
	if err != nil {
		return err
	}

	for _, item := range input.Items {
		// TODO: find product in the repository database here
		fakeProductName := "Product " + item.ProductId
		fakeProductPrice := 10.50

		// create fake order item
		i := entity.NewOrderItemEntity(fakeProductName, fakeProductPrice, item.Qtd)

		// add items to order
		order.AddItem(i)
	}

	// TODO: save the order in the repository database here

	var eventItems []event.OrderItem
	for _, item := range order.GetItems() {
		eventItems = append(eventItems, event.OrderItem{
			ProductName: item.GetProductName(),
			TotalPrice:  item.GetTotalPrice(),
			Quantity:    item.GetQuantity(),
		})
	}

	// publish event OrderCreatedEvent passing the order data
	err = u.publisher.Publish(ctx, event.OrderCreatedEvent{
		Id:         order.GetID(),
		TotalPrice: order.GetTotalPrice(),
		Status:     order.GetStatus(),
		Items:      eventItems,
	})
	if err != nil {
		return err
	}
	return nil
}
