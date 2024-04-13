package usecase

import (
	"fmt"

	"github.com/briannbig/event-driven/internal/domain/event"
)

type SendOrderEmailUseCase struct {
}

func NewSendOrderEmailUseCase() *SendOrderEmailUseCase {
	return &SendOrderEmailUseCase{}
}

func (h *SendOrderEmailUseCase) Execute(payload *event.OrderCreatedEvent) error {
	fmt.Println("--- SendOrderEmailUseCase ---")
	fmt.Printf("--- MAIL Order Created: R$ %f \n", payload.TotalPrice)
	return nil
}
