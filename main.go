package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	fmt.Println("hello")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/customer/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	e.GET("/customer/:id/instagram", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	e.POST("/customer/:id/facebook_token", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Transaction interface {
	Commit() error
	Rollback()
	GetTx() *gorm.DB
}

type Tx struct {
	tx *gorm.DB
}

func (t Tx) Commit() error {
	return t.tx.Commit().Error
}

func (t Tx) Rollback() {
	t.tx.Rollback()
}

func (t Tx) GetTx() *gorm.DB {
	return t.tx
}

type Filter interface {
	GenerateMods(db *gorm.DB) *gorm.DB
}

type DbClient interface {
	BeginTransaction() Transaction

	First(model interface{}, filter Filter) error
}

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type baseRepository struct {
	dbClient DbClient
}

type BaseRepository interface {
	BeginTransaction() Transaction
}

func (b *baseRepository) BeginTransaction() Transaction {
	return b.dbClient.BeginTransaction()
}

type customerRepository struct {
	dbClient DbClient
}

func (c *customerRepository) FindByID(id string) (*Customer, error) {
	c.dbClient.First()
}

type CustomerRepository interface {
	FindByID(id string) (*Customer, error)
}

type customerService struct {
	baseRepository     BaseRepository
	customerRepository CustomerRepository
}

func (s *customerService) GetCustomer(ctx context.Context, id string) (*Customer, error) {

}

type CustomerService interface {
	GetCustomer(c context.Context, id string) (*Customer, error)
}

type Presenter interface {
	Generate(error, any) (int, interface{})
}

type presenter struct {
}

func (p *presenter) Generate(err error, body any) (int, interface{}) {
	return http.StatusNotImplemented, nil
}

type customerController struct {
	customerService CustomerService
	presenter       Presenter
}

func (ctr *customerController) GetCustomer(c echo.Context) error {
	customer, err := ctr.customerService.GetCustomer(context.Background(), c.Param("id"))
	return c.JSON(ctr.presenter.Generate(err, customer))
}
