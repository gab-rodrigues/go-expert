package usecase

import (
	"clean-arch/internal/entity"
)

type ListOrderOutputDTO struct {
	Orders []OrderOutputDTO `json:"orders"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() (ListOrderOutputDTO, error) {
	orders, err := l.OrderRepository.ListAll()
	if err != nil {
		return ListOrderOutputDTO{}, err
	}

	var ordersDTO []OrderOutputDTO
	for _, order := range orders {
		ordersDTO = append(ordersDTO, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return ListOrderOutputDTO{Orders: ordersDTO}, nil
}
