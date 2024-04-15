package usecase

import "github.com/devfullcycle/20-CleanArch/internal/entity"

type ListOrderOutputDTO []OrderOutputDTO

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (uc *ListOrderUseCase) Execute() (ListOrderOutputDTO, error) {
	orders, err := uc.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	ordersDTO := make(ListOrderOutputDTO, len(orders))
	for i, order := range orders {
		ordersDTO[i] = OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}

	return ordersDTO, nil
}
