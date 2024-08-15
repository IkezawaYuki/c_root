package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DatabaseUser,
		config.Env.DatabasePass,
		config.Env.DatabaseHost,
		config.Env.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr: config.Env.RedisAddr,
		DB:   0,
	})

	customerController := NewCustomerController(db)
	adminController := NewAdminController(db)
	authService := NewAuthService(db, redisDb)
	customerAuthMiddleware := NewCustomerAuthMiddleware(authService)
	adminAuthMiddleware := NewAdminAuthMiddleware(authService)
	badgeAuthMiddleware := NewBadgeAuthMiddleware(authService)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/customer/login", customerController.Login)

	customerHandler := e.Group("/customer")
	customerHandler.Use(customerAuthMiddleware)
	customerHandler.GET("/:id", func(c echo.Context) error {
		return customerController.GetCustomer(c)
	})
	customerHandler.GET("/:id/instagram", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	customerHandler.POST("/:id/facebook_token", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	adminHandler := e.Group("/admin")
	adminHandler.Use(adminAuthMiddleware)
	adminHandler.GET("/:id", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, nil)
	})
	adminHandler.POST("/register/customer", adminController.RegisterCustomer)

	badgeHandler := e.Group("/badge")
	badgeHandler.Use(badgeAuthMiddleware)
	badgeHandler.GET("/:id", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, nil)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CustomerDto struct {
	ID             string         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	Email          string         `gorm:"column:email"`
	Password       string         `gorm:"column:password"`
	FacebookToken  sql.NullString `gorm:"column:facebook_token"`
	StartDate      sql.NullTime   `gorm:"column:start_date"`
	InstagramID    sql.NullString `gorm:"column:instagram_id"`
	InstagramName  sql.NullString `gorm:"column:instagram_name"`
	DeleteHashFlag int            `gorm:"column:delete_hash_flag"`
	gorm.Model
}

func (c *CustomerDto) TableName() string {
	return "customers"
}

func (c *CustomerDto) ConvertToCustomer() *Customer {
	return &Customer{
		ID:       c.ID,
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
	}
}

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := r.client.Set(ctx, key, value, expiration).Result()
	return err
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
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

func (c *CustomerRepository) FindByEmail(ctx context.Context, email string) (*Customer, error) {
	var customer CustomerDto
	result := c.db.WithContext(ctx).First(&customer, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
	}
	return customer.ConvertToCustomer(), nil
}

func (c *CustomerRepository) Save(ctx context.Context, customer *CustomerDto) *gorm.DB {
	return c.db.WithContext(ctx).Save(customer)
}

type CustomerService struct {
	baseRepository     BaseRepository
	customerRepository CustomerRepository
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	return s.customerRepository.FindByID(ctx, id)
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	return s.customerRepository.FindByEmail(ctx, email)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *Customer) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.customerRepository.Save(ctx, &CustomerDto{
		ID:             uuid.New().String(),
		Name:           customer.Name,
		Email:          customer.Email,
		Password:       string(passwordHash),
		FacebookToken:  sql.NullString{},
		StartDate:      sql.NullTime{},
		InstagramID:    sql.NullString{},
		InstagramName:  sql.NullString{},
		DeleteHashFlag: 0,
	}).Error
}

type Presenter struct {
}

func (p *Presenter) Generate(err error, body any) (int, interface{}) {
	if err != nil {
		return http.StatusOK, body
	}
	return http.StatusNotImplemented, nil
}

type CustomerController struct {
	customerUsecase CustomerUsecase
	presenter       Presenter
}

type CustomerUsecase struct {
	customerService CustomerService
	authService     AuthService
}

func (c *CustomerUsecase) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	return c.customerService.GetCustomer(ctx, id)
}

func (c *CustomerUsecase) Login(ctx context.Context, user *User) (string, error) {
	customer, err := c.customerService.GetCustomerByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return c.authService.GenerateJWTCustomer(customer)
}

func NewCustomerController(db *gorm.DB) CustomerController {
	return CustomerController{
		customerUsecase: CustomerUsecase{
			customerService: CustomerService{
				baseRepository: BaseRepository{
					db: db,
				},
				customerRepository: CustomerRepository{
					db: db,
				},
			},
		},
		presenter: Presenter{},
	}
}

type AdminController struct {
	adminUsecase AdminUsecase
	presenter    Presenter
}

type AdminUsecase struct {
	adminService    AdminService
	authService     AuthService
	customerService CustomerService
}

func (a *AdminUsecase) RegisterCustomer(ctx context.Context, customer *Customer) error {
	return a.customerService.CreateCustomer(ctx, customer)
}

func (a *AdminController) RegisterCustomer(c echo.Context) error {
	var customer Customer
	if err := c.Bind(&customer); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := a.adminUsecase.RegisterCustomer(c.Request().Context(), &customer)
	return c.JSON(a.presenter.Generate(err, customer))
}

func NewAdminController(db *gorm.DB) AdminController {
	return AdminController{
		adminUsecase: AdminUsecase{
			adminService: AdminService{
				customerRepository: CustomerRepository{db: db},
				adminRepository:    AdminRepository{db: db},
			},
		},
	}
}

type AdminService struct {
	customerRepository CustomerRepository
	adminRepository    AdminRepository
}

type AdminRepository struct {
	db *gorm.DB
}

type AuthController struct {
	customerService CustomerService
	authService     AuthService
}

//func (a *AuthController) LoginCustomer(ctx context.Context, email string, password string) (*Customer, error) {
//	customer, err := a.customerService.GetCustomerByEmail(ctx, email)
//	if err != nil {
//		return nil, err
//	}
//	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password)); err != nil {
//		return nil, err
//	}
//	return nil, nil
//}

func NewAuthController(db *gorm.DB, redisCli *redis.Client) AuthController {
	return AuthController{
		authService: AuthService{
			customerRepository: CustomerRepository{
				db: db,
			},
			redisClient: RedisClient{
				client: redisCli,
			},
		},
		customerService: CustomerService{
			baseRepository: BaseRepository{
				db: db,
			},
			customerRepository: CustomerRepository{
				db: db,
			},
		},
	}
}

func NewAuthService(db *gorm.DB, client *redis.Client) AuthService {
	return AuthService{
		customerRepository: CustomerRepository{
			db: db,
		},
		redisClient: RedisClient{
			client: client,
		},
	}
}

type AuthService struct {
	customerRepository CustomerRepository
	redisClient        RedisClient
}

func (a *AuthService) IsCustomerIsLogin(token string) (bool, error) {
	return true, nil
}

func (a *AuthService) IsAdminLogin() (bool, error) {
	return true, nil
}

func (a *AuthService) GenerateJWTCustomer(c *Customer) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "c_root",
		"aud":   "customer",
		"sub":   c.ID,
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) CheckPassword(user *User, customer *Customer) error {
	return bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(user.Password))
}

func NewBadgeAuthMiddleware(srv AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func NewAdminAuthMiddleware(srv AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("middleware")
			login, err := srv.IsAdminLogin()
			if !login {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
			return next(c)
		}
	}
}

func NewCustomerAuthMiddleware(srv AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("middleware")
			token := c.Request().Header.Get("Authorization")
			slog.Info(token)
			login, err := srv.IsCustomerIsLogin(token)
			if !login {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
			return next(c)
		}
	}
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ctr *CustomerController) Login(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	customerId := c.Param("id")
	ctx := c.Request().Context()
	customer, err := ctr.customerUsecase.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}
