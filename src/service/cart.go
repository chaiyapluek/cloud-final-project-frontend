package service

import (
	"errors"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
	"dev.chaiyapluek.cloud.final.frontend/src/repository"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/cart"
)

type CartService interface {
	GetUserCart(userId, locationId string) (*cart.CartProps, error)
	AddCartItem(userId, locationId string, item *entity.AddCartItemRequest) (*entity.AddCartItem, error)
	DeleteCartItem(cartId string, itemId int) (int, int, error)
	Checkout(cartId, userId, address string) error
}

type cartService struct {
	cartRepo repository.CartPepository
}

func NewCartService(cartRepo repository.CartPepository) CartService {
	return &cartService{
		cartRepo: cartRepo,
	}
}

func (s *cartService) GetUserCart(userId, locationId string) (*cart.CartProps, error) {
	userCart, err := s.cartRepo.GetUserCart(userId, locationId)
	if err != nil {
		return nil, err
	}

	totalPrice := 0
	cartItems := []*cart.CartItem{}
	for _, item := range userCart.CartItems {
		itemSteps := []*cart.CartItemStep{}
		for _, step := range item.Steps {
			itemSteps = append(itemSteps, &cart.CartItemStep{
				Step:    step.Step,
				Options: step.Options,
			})
		}
		cartItems = append(cartItems, &cart.CartItem{
			ItemId:     item.ItemId,
			MenuId:     item.MenuId,
			MenuName:   item.MenuName,
			TotalPrice: item.TotalPrice,
			Quantity:   item.Quantity,
			Steps:      itemSteps,
		})
		totalPrice += item.TotalPrice * item.Quantity
	}

	return &cart.CartProps{
		CartId:       userCart.CartId,
		LocationId:   userCart.LocationId,
		LocationName: userCart.LocationName,
		TotalPrice:   totalPrice,
		CartItems:    cartItems,
	}, nil
}

func (s *cartService) AddCartItem(userId, locationId string, item *entity.AddCartItemRequest) (*entity.AddCartItem, error) {
	userCart, err := s.cartRepo.GetUserCart(userId, locationId)
	if err != nil {
		return nil, err
	}
	if userCart == nil {
		return nil, errors.New("cart not found")
	}

	resp, err := s.cartRepo.AddCartItem(userCart.CartId, item)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *cartService) DeleteCartItem(cartId string, itemId int) (int, int, error) {
	resp, err := s.cartRepo.DeleteCartItem(cartId, itemId)
	if err != nil {
		return 0, 0, err
	}

	return resp.TotalPrice, resp.TotalItem, nil
}

func (s *cartService) Checkout(cartId, userId, address string) error {
	req := &entity.CheckoutRequest{
		CartId:  cartId,
		UserId:  userId,
		Address: address,
	}

	err := s.cartRepo.Checkout(req)
	if err != nil {
		return err
	}

	return nil
}
