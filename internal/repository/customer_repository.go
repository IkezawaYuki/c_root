package repository

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (c *CustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	var customer domain.CustomerDto
	result := c.db.WithContext(ctx).First(&customer, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, result.Error
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) FindByIDTx(ctx context.Context, id string, tx *gorm.DB) (*domain.Customer, error) {
	var customer domain.CustomerDto
	result := tx.WithContext(ctx).First(&customer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, result.Error
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var customer domain.CustomerDto
	result := c.db.WithContext(ctx).First(&customer, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) Save(ctx context.Context, customer *domain.CustomerDto) *gorm.DB {
	return c.db.WithContext(ctx).Save(customer)
}
