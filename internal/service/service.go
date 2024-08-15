package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/infrastructure"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
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
	customer.ID = uuid.New().String()
	if err := s.customerRepository.Save(ctx, &domain.CustomerDto{
		ID:             customer.ID,
		Name:           customer.Name,
		Email:          customer.Email,
		Password:       string(passwordHash),
		DeleteHashFlag: 0,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEmail
		}
		return err
	}
	return nil
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
	slog.Info("IsCustomerIsLogin is invoked")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		slog.Info(err.Error())
		return "", domain.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", domain.ErrAuthorization
	}
	if !claims.VerifyAudience("customer", true) {
		return "", domain.ErrAuthentication
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

func (a *AdminService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return a.customerRepository.FindByID(ctx, id)
}

type InstagramService struct {
	httpClient infrastructure.HttpClient
}

func NewInstagramService(httpClient infrastructure.HttpClient) InstagramService {
	return InstagramService{
		httpClient: httpClient,
	}
}

const getMediaList = "/media"

func (i *InstagramService) GetMediaList(ctx context.Context, facebookToken string) ([]string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		getMediaList,
		nil,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaList
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return detail.ConvertToInstagramMediaList(), nil
}

const getMediaDetail = "/media/detail"

func (i *InstagramService) GetMediaDetail(ctx context.Context, facebookToken string, id string) (*domain.InstagramMediaDetail, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		getMediaDetail,
		nil,
		fmt.Sprintf("Bearer %s", facebookToken),
	)
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaDetail
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

type WordpressService struct {
	httpClient infrastructure.HttpClient
}

func (w *WordpressService) Post() {

}

func (w *WordpressService) UploadFile() {

}

type MetaService struct {
	httpClient infrastructure.HttpClient
}

func (m *MetaService) GetLongToken() {

}
