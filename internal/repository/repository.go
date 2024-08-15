package repository

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return CustomerRepository{db: db}
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

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}

func (b *BaseRepository) Begin() *gorm.DB {
	return b.db.Begin()
}

func (b *BaseRepository) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (b *BaseRepository) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return AdminRepository{db: db}
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) RedisClient {
	return RedisClient{client: client}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := r.client.Set(ctx, key, value, expiration).Result()
	return err
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
