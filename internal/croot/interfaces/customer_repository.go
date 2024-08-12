package interfaces

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/croot/domain"
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre"
)

type CustomerRepository interface {
	FindByIDWithTX(ctx context.Context, id string, tx infrastructre.Tx) (*domain.Customer, error)
	FindAllWithTx(ctx context.Context, page int, tx infrastructre.Tx) []domain.Customer
	CreateWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer
	UpdateWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer
	DeleteWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer
	FindByID(ctx context.Context, id string) (*domain.Customer, error)
	FindAll(ctx context.Context, page int) []domain.Customer
	Create(ctx context.Context, customer domain.Customer) domain.Customer
	Update(ctx context.Context, customer domain.Customer) domain.Customer
	Delete(ctx context.Context, customer domain.Customer) domain.Customer
}

type customerRepository struct {
	client infrastructre.DbClient
}

func NewCustomerAdapter() CustomerRepository {
	return &customerRepository{}
}

func (c *customerRepository) FindByIDWithTX(ctx context.Context, id string, tx infrastructre.Tx) (*domain.Customer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) FindAllWithTx(ctx context.Context, page int, tx infrastructre.Tx) []domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) CreateWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) UpdateWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) DeleteWithTx(ctx context.Context, customer domain.Customer, tx infrastructre.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) FindAll(ctx context.Context, page int) []domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) Create(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) Update(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerRepository) Delete(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}

const findCustomerByID = `select id, name, email, password, wordpress_url, facebook_token, start_date, instagram_id, instagram_name from customers where id = ?;`

func (c *customerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	row := c.client.QueryRowContext(ctx, findCustomerByID, id)
	var customer domain.Customer
	if err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Password,
		&customer.WordpressURL,
		&customer.FacebookToken,
		&customer.StartDate,
		&customer.InstagramID,
		&customer.InstagramName,
		&customer.DeleteHashFlag,
	); err != nil {
		return nil, err
	}
	return &customer, nil
}
