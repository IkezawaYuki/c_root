package service

import (
	"context"
	"github.com/IkezawaYuki/c_root/domain/entity"
)

type CustomerService interface {
	FindByID(ctx context.Context, customerID string) entity.Customer
	FindAll(ctx context.Context, page int) []entity.Customer
	Create(ctx context.Context, customer entity.Customer) entity.Customer
	Update(ctx context.Context, customer entity.Customer) entity.Customer
	Delete(ctx context.Context, customer entity.Customer) entity.Customer
}
