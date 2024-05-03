package orders

import (
	"fmt"

	"github.com/andrian0vv/test-go-project/internal/api"
)

type OutOfStockError struct {
	productID int64
}

func (e *OutOfStockError) Error() string {
	return fmt.Sprintf("product %d is out of stock", e.productID)
}

func (e *OutOfStockError) Code() string {
	return api.CodeOutOfStock
}

func newOutOfStockError(productID int64) *OutOfStockError {
	return &OutOfStockError{
		productID: productID,
	}
}

type ProductNotFoundError struct {
	productID int64
}

func (e *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product %d not found", e.productID)
}

func (e *ProductNotFoundError) Code() string {
	return api.CodeProductNotFound
}

func newProductNotFoundError(productID int64) *ProductNotFoundError {
	return &ProductNotFoundError{
		productID: productID,
	}
}
