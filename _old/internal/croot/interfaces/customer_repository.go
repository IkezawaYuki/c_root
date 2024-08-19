package interfaces

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/croot/domain"
	"github.com/IkezawaYuki/popple/internal/croot/infrastructre"
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
	Update(ctx context.Context, customer domain.Customer) error
	Delete(ctx context.Context, customer domain.Customer) error
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

const findAllCustomerQuery = `select id, name, email, password, wordpress_url, facebook_token, start_date, instagram_id, instagram_name, delete_hash_flag
from customers
limit ? offset ?;
`

func (c *customerRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Customer, error) {
	rows, err := c.client.QueryContext(ctx, findAllCustomerQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var customers []domain.Customer
	for rows.Next() {
		var customer domain.Customer
		if  rows.Scan();
	}
}

const insertCustomerQuery = `
insert into customers (name, email, password, wordpress_url, facebook_token, start_date, instagram_id, instagram_name, delete_hash_flag)
values(?, ?, ?, ?, ?, ?, ?, ?, ?);
`

func (c *customerRepository) Create(ctx context.Context, customer domain.Customer) (*domain.Customer, error) {
	result, err := c.client.ExecContext(ctx, insertCustomerQuery,
		customer.Name,
		customer.Email,
		customer.Password,
		customer.WordpressURL,
		customer.FacebookToken,
		customer.StartDate,
		customer.InstagramID,
		customer.InstagramName,
		customer.DeleteHashFlag,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	customer.ID = int(id)
	return &customer, nil
}

const customerUpdateQuery = `update customers set 
name = ?,
email = ?,
password = ?,
wordpress_url = ?,
start_date = ?,
instagram_id = ?,
instagram_name = ?,
delete_hash_flag = ?
where id = ?;`

func (c *customerRepository) Update(ctx context.Context, customer domain.Customer) error {
	_, err := c.client.ExecContext(ctx, customerUpdateQuery,
		customer.Name,
		customer.Email,
		customer.Password,
		customer.WordpressURL,
		customer.StartDate,
		customer.InstagramID,
		customer.InstagramName,
		customer.DeleteHashFlag,
		customer.ID,
	)
	return err
}

const customerDeleteQuery = "delete from customer where id = ?"

func (c *customerRepository) Delete(ctx context.Context, customer domain.Customer) error {
	_, err := c.client.ExecContext(ctx, customerDeleteQuery, customer.ID)
	return err
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
