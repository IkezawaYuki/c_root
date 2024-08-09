package interfaces

import (
	"context"
	"database/sql"
	"github.com/IkezawaYuki/c_root/internal/croot/domain"
)

type customerAdapter struct {
	db *sql.DB
}

func NewCustomerAdapter() domain.CustomerRepository {
	return &customerAdapter{}
}

func (c *customerAdapter) FindByIDWithTX(ctx context.Context, id string, tx domain.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) FindAllWithTx(ctx context.Context, page int, tx domain.Tx) []domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) CreateWithTx(ctx context.Context, customer domain.Customer, tx domain.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) UpdateWithTx(ctx context.Context, customer domain.Customer, tx domain.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) DeleteWithTx(ctx context.Context, customer domain.Customer, tx domain.Tx) domain.Customer {
	//TODO implement me
	panic("implement me")
}

const findCustomerByID = `select id, name, email, password, wordpress_url, facebook_token, start_date, instagram_id, instagram_name from customers where id = ?;`

func (c *customerAdapter) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	stmt, err := c.db.PrepareContext(ctx, findCustomerByID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	row := stmt.QueryRowContext(ctx, id)
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

func (c *customerAdapter) FindAll(ctx context.Context, page int) []domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) Create(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) Update(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}

func (c *customerAdapter) Delete(ctx context.Context, customer domain.Customer) domain.Customer {
	//TODO implement me
	panic("implement me")
}
