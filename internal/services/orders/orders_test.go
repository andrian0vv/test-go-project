package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/andrian0vv/test-go-project/internal/models"
	"github.com/andrian0vv/test-go-project/internal/services/orders/mocks"
)

func TestService_CreateOrder(t *testing.T) {
	testCases := []struct {
		name                string
		dbFn                func(*gomock.Controller) *mocks.Mockdatabase
		orderRepositoryFn   func(*gomock.Controller) *mocks.MockordersRepository
		productRepositoryFn func(*gomock.Controller) *mocks.MockproductsRepository
		in                  Order
		wantErr             bool
	}{
		{
			name: "get products error",
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1, 2}).
					Return(nil, errors.New("some error"))
				return r
			},
			in: Order{
				Products: []Product{
					{ProductID: 1},
					{ProductID: 2},
				},
			},
			wantErr: true,
		},
		{
			name: "out of stock",
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1}).
					Return(map[int64]models.Product{1: {Quantity: 9}}, nil)
				return r
			},
			in: Order{
				Products: []Product{
					{ProductID: 1, Quantity: 10},
				},
			},
			wantErr: true,
		},
		{
			name: "product not found",
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1}).
					Return(map[int64]models.Product{}, nil)
				return r
			},
			in: Order{
				Products: []Product{
					{ProductID: 1, Quantity: 10},
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			dbFn: func(ctrl *gomock.Controller) *mocks.Mockdatabase {
				d := mocks.NewMockdatabase(ctrl)
				d.EXPECT().Tx(gomock.Any(), gomock.Any()).
					Do(func(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) {
						assert.NoError(t, fn(ctx, nil))
					}).
					Return(nil)
				return d
			},
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1}).
					Return(map[int64]models.Product{
						1: {
							Quantity: 100,
							Price:    5,
						},
					}, nil)
				r.EXPECT().
					WriteOff(gomock.Any(), gomock.Any(), map[int64]int{1: 2}).
					Return(nil)
				return r
			},
			orderRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockordersRepository {
				r := mocks.NewMockordersRepository(ctrl)
				r.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any(), models.Order{
						UserID: 111,
						Products: []models.OrderProduct{{
							ProductID: 1,
							Quantity:  2,
							Price:     5 * 2,
						}},
					}).
					Return(nil)
				return r
			},
			in: Order{
				UserID: 111,
				Products: []Product{
					{ProductID: 1, Quantity: 2},
				},
			},
		},
		{
			name: "create order error",
			dbFn: func(ctrl *gomock.Controller) *mocks.Mockdatabase {
				d := mocks.NewMockdatabase(ctrl)
				d.EXPECT().Tx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
						err := fn(ctx, nil)
						assert.Error(t, err)
						return err
					})
				return d
			},
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1}).
					Return(map[int64]models.Product{
						1: {
							Quantity: 100,
							Price:    5,
						},
					}, nil)
				return r
			},
			orderRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockordersRepository {
				r := mocks.NewMockordersRepository(ctrl)
				r.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any(), models.Order{
						UserID: 111,
						Products: []models.OrderProduct{{
							ProductID: 1,
							Quantity:  2,
							Price:     5 * 2,
						}},
					}).
					Return(errors.New("some error"))
				return r
			},
			in: Order{
				UserID: 111,
				Products: []Product{
					{ProductID: 1, Quantity: 2},
				},
			},
			wantErr: true,
		},
		{
			name: "write off error",
			dbFn: func(ctrl *gomock.Controller) *mocks.Mockdatabase {
				d := mocks.NewMockdatabase(ctrl)
				d.EXPECT().Tx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
						err := fn(ctx, nil)
						assert.Error(t, err)
						return err
					})
				return d
			},
			productRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockproductsRepository {
				r := mocks.NewMockproductsRepository(ctrl)
				r.EXPECT().
					GetProducts(gomock.Any(), []int64{1}).
					Return(map[int64]models.Product{
						1: {
							Quantity: 100,
							Price:    5,
						},
					}, nil)
				r.EXPECT().
					WriteOff(gomock.Any(), gomock.Any(), map[int64]int{1: 2}).
					Return(errors.New("some error"))
				return r
			},
			orderRepositoryFn: func(ctrl *gomock.Controller) *mocks.MockordersRepository {
				r := mocks.NewMockordersRepository(ctrl)
				r.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any(), models.Order{
						UserID: 111,
						Products: []models.OrderProduct{{
							ProductID: 1,
							Quantity:  2,
							Price:     5 * 2,
						}},
					}).
					Return(nil)
				return r
			},
			in: Order{
				UserID: 111,
				Products: []Product{
					{ProductID: 1, Quantity: 2},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var db *mocks.Mockdatabase
			if tc.dbFn != nil {
				db = tc.dbFn(ctrl)
			}

			var or *mocks.MockordersRepository
			if tc.orderRepositoryFn != nil {
				or = tc.orderRepositoryFn(ctrl)
			}

			var pr *mocks.MockproductsRepository
			if tc.productRepositoryFn != nil {
				pr = tc.productRepositoryFn(ctrl)
			}

			s := New(db, or, pr)

			err := s.CreateOrder(context.Background(), tc.in)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
