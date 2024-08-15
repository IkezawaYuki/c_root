package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type CustomerService struct {
	baseRepository     repository.BaseRepository
	customerRepository repository.CustomerRepository
}

func NewCustomerService(baseRepo repository.BaseRepository, customerRepo repository.CustomerRepository) CustomerService {
	return CustomerService{
		baseRepository:     baseRepo,
		customerRepository: customerRepo,
	}
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.customerRepository.FindByID(ctx, id)
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	return s.customerRepository.FindByEmail(ctx, email)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.customerRepository.Save(ctx, &domain.CustomerDto{
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

type AuthService struct {
	customerRepository repository.CustomerRepository
	redisClient        repository.RedisClient
}

func NewAuthService(customerRepo repository.CustomerRepository, redisClient repository.RedisClient) AuthService {
	return AuthService{
		customerRepository: customerRepo,
		redisClient:        redisClient,
	}
}

func (a *AuthService) IsCustomerIsLogin(tokenString string) (string, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	if !claims.VerifyAudience("customer", true) {
		return "", fmt.Errorf("not authentication")
	}
	return claims["sub"].(string), nil
}

func (a *AuthService) IsAdminLogin() (bool, error) {
	return true, nil
}

func (a *AuthService) GenerateJWTCustomer(c *domain.Customer) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "customer",
		"sub":   c.ID,
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) CheckPassword(user *domain.User, customer *domain.Customer) error {
	return bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(user.Password))
}

type AdminService struct {
	customerRepository repository.CustomerRepository
	adminRepository    repository.AdminRepository
}

func NewAdminService(customerRepo repository.CustomerRepository, adminRepo repository.AdminRepository) AdminService {
	return AdminService{
		customerRepository: customerRepo,
		adminRepository:    adminRepo,
	}
}
