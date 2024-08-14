package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	fmt.Println("hello")

	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	customerController := NewCustomerController(db)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/customer/:id", func(c echo.Context) error {
		return customerController.GetCustomer(c)
	})
	e.GET("/customer/:id/instagram", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	e.POST("/customer/:id/facebook_token", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CustomerDto struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	EMail          string         `gorm:"column:email"`
	Password       string         `gorm:"column:password"`
	FacebookToken  sql.NullString `gorm:"column:facebook_token"`
	StartDate      sql.NullTime   `gorm:"column:start_date"`
	InstagramID    sql.NullString `gorm:"column:instagram_id"`
	InstagramName  sql.NullString `gorm:"column:instagram_name"`
	DeleteHashFlag int            `gorm:"column:delete_hash_flag"`
	gorm.Model
}

func (c *CustomerDto) ConvertToCustomer() *Customer {
	return &Customer{
		ID:   c.ID,
		Name: c.Name,
	}
}

type BaseRepository struct {
	db *gorm.DB
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

type CustomerRepository struct {
	db *gorm.DB
}

func (c *CustomerRepository) FindByID(ctx context.Context, id string) (*Customer, error) {
	var customer CustomerDto
	result := c.db.WithContext(ctx).First(&customer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) FindByIDTx(ctx context.Context, id string, tx *gorm.DB) (*Customer, error) {
	var customer CustomerDto
	result := tx.WithContext(ctx).First(&customer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) SaveTx(customer *Customer, tx *gorm.DB) *gorm.DB {
	return tx.Save(customer)
}

type CustomerService struct {
	baseRepository     BaseRepository
	customerRepository CustomerRepository
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	customer, err := s.customerRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

type Presenter struct {
}

func (p *Presenter) Generate(err error, body any) (int, interface{}) {
	return http.StatusNotImplemented, nil
}

type CustomerController struct {
	customerService CustomerService
	presenter       Presenter
}

func NewCustomerController(db *gorm.DB) CustomerController {
	return CustomerController{
		customerService: CustomerService{
			baseRepository: BaseRepository{
				db: db,
			},
			customerRepository: CustomerRepository{
				db: db,
			},
		},
		presenter: Presenter{},
	}
}

func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	customerId := c.Param("id")
	ctx := c.Request().Context()
	customer, err := ctr.customerService.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}
