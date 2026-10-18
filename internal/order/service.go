package order

import (
	"context"
	"e-commerce/internal/cart"
	"e-commerce/internal/product"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo           *Repository
	cartService    *cart.Service
	productService *product.Service
}

func NewService(repo *Repository, cartService *cart.Service, productService *product.Service) *Service {
	return &Service{repo: repo, cartService: cartService, productService: productService}
}

func (s *Service) Create(ctx context.Context, uid uuid.UUID) error {
	cartItems, err := s.cartService.GetItemsByUserID(ctx, uid)
	if err != nil {
		return err
	}

	if len(cartItems) == 0 {
		return errors.New("cart kosong")
	}

	order := &Order{
		UserID: uid,
		Status: "pending",
	}

	for _, cartItem := range cartItems {
		productID := cartItem.ProductID.String()

		p, err := s.productService.GetById(ctx, productID)
		if err != nil {
			return err
		}

		// Perika apakah stok mencukupi
		if p.Stock < cartItem.Quantity {
			return errors.New("Stok product tidak mencukupi")
		}

		subtotal := p.Price * float64(cartItem.Quantity)

		item := OrderItem{
			ProductID:   cartItem.ProductID,
			ProductName: p.Name,
			UnitPrice:   p.Price,
			Quantity:    cartItem.Quantity,
			Subtotal:    subtotal,
		}

		order.Items = append(order.Items, item)
		order.Total += subtotal
	}

	return s.repo.Create(ctx, order)
}
