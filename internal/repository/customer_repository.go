package repository

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (c *CustomerRepository) FindAll(ctx context.Context) ([]model.Customer, error) {
	var dto []model.Customer
	err := c.db.WithContext(ctx).Find(&dto).Error
	return dto, err
}

func (c *CustomerRepository) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.WithContext(ctx).First(&customer, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, objects.ErrNotFound
		}
		return nil, result.Error
	}
	return &customer, nil
}

func (c *CustomerRepository) FindByIDTx(ctx context.Context, id string, tx *gorm.DB) (*model.Customer, error) {
	var customer model.Customer
	result := tx.WithContext(ctx).First(&customer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, result.Error
	}
	return &customer, nil
}

func (c *CustomerRepository) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.WithContext(ctx).First(&customer, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, objects.ErrNotFound
		}
		return nil, result.Error
	}
	return &customer, nil
}

func (c *CustomerRepository) Save(ctx context.Context, customer *model.Customer) *gorm.DB {
	return c.db.WithContext(ctx).Save(customer)
}
