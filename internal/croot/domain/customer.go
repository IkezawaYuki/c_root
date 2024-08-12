package domain

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre"
	"time"
)

type Customer struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	WordpressURL   string     `json:"wordpress_url"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag int        `json:"delete_hash_flag"`
}

type CustomerRepository interface {
	FindByIDWithTX(ctx context.Context, id string, tx infrastructre.Tx) (Customer, error)
	FindAllWithTx(ctx context.Context, page int, tx Tx) []Customer
	CreateWithTx(ctx context.Context, customer Customer, tx Tx) Customer
	UpdateWithTx(ctx context.Context, customer Customer, tx Tx) Customer
	DeleteWithTx(ctx context.Context, customer Customer, tx Tx) Customer

	FindByID(ctx context.Context, id string) (*Customer, error)
	FindAll(ctx context.Context, page int) []Customer
	Create(ctx context.Context, customer Customer) Customer
	Update(ctx context.Context, customer Customer) Customer
	Delete(ctx context.Context, customer Customer) Customer
}
