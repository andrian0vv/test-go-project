//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks

package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/andrian0vv/test-go-project/internal/models"
)

type database interface {
	Tx(context.Context, func(context.Context, *sqlx.Tx) error) error
}

type ordersRepository interface {
	CreateOrder(context.Context, *sqlx.Tx, models.Order) error
}

type productsRepository interface {
	GetProducts(ctx context.Context, ids []int64) (map[int64]models.Product, error)
	WriteOff(ctx context.Context, tx *sqlx.Tx, products map[int64]int) error
}

type Service struct {
	database           database
	ordersRepository   ordersRepository
	productsRepository productsRepository
}

func New(database database, ordersRepository ordersRepository, productsRepository productsRepository) *Service {
	return &Service{
		database:           database,
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order Order) error {
	products, err := s.productsRepository.GetProducts(ctx, order.productIDs())
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	model := models.Order{
		UserID: order.UserID,
	}

	requiredProducts := requiredProducts(order)

	var validateionErrors []error
	for productID, quantity := range requiredProducts {
		product, exists := products[productID]
		if !exists {
			validateionErrors = append(validateionErrors, newProductNotFoundError(productID))
			continue
		}

		if product.Quantity < quantity {
			validateionErrors = append(validateionErrors, newOutOfStockError(productID))
			continue
		}

		model.Products = append(model.Products, models.OrderProduct{
			ProductID: productID,
			Price:     product.Price * int64(quantity),
			Quantity:  quantity,
		})
	}

	if len(validateionErrors) > 0 {
		return errors.Join(validateionErrors...)
	}

	err = s.database.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err = s.ordersRepository.CreateOrder(ctx, tx, model)
		if err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		err = s.productsRepository.WriteOff(ctx, tx, requiredProducts)
		if err != nil {
			// todo explicitly check constraint 'products_quantity_check'
			return fmt.Errorf("write off products: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("tx: %w", err)
	}

	return nil
}

func requiredProducts(order Order) map[int64]int {
	result := make(map[int64]int)
	for _, p := range order.Products {
		result[p.ProductID] += p.Quantity
	}

	return result
}
