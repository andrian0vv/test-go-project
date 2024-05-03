package orders

type Order struct {
	UserID   int64
	Products []Product
}

type Product struct {
	ProductID int64
	Quantity  int
}

func (o Order) productIDs() []int64 {
	result := make([]int64, 0, len(o.Products))
	for _, p := range o.Products {
		result = append(result, p.ProductID)
	}

	return result
}
