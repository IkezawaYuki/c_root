package service

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strings"
	"time"
)

type AuthService struct {
	customerRepository *repository.CustomerRepository
	redisClient        *repository.RedisClient
}

func NewAuthService(customerRepo *repository.CustomerRepository, redisClient *repository.RedisClient) *AuthService {
	return &AuthService{
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

func (a *AuthService) IsAdminLogin(tokenString string) (string, error) {
	slog.Info("IsAdminLogin is invoked")
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
	if !claims.VerifyAudience("admin", true) {
		return "", domain.ErrAuthentication
	}
	return claims["sub"].(string), nil
}

func (a *AuthService) GenerateJWTCustomer(c *domain.Customer) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "customer",
		"sub":   c.UUID,
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) GenerateJWTAdmin(admin *domain.Admin) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "admin",
		"sub":   admin.UUID,
		"email": admin.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) CheckPassword(user *domain.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return fmt.Errorf("password is incorrect: %s, %v", err.Error(), domain.ErrAuthorization)
	}
	return nil
}

type AdminService struct {
	customerRepository *repository.CustomerRepository
	adminRepository    *repository.AdminRepository
}

func NewAdminService(customerRepo *repository.CustomerRepository, adminRepo *repository.AdminRepository) *AdminService {
	return &AdminService{
		customerRepository: customerRepo,
		adminRepository:    adminRepo,
	}
}

func (a *AdminService) GetCustomerByID(ctx context.Context, id string) (*domain.Customer, error) {
	return a.customerRepository.FindByID(ctx, id)
}

func (a *AdminService) FindAll(ctx context.Context) ([]domain.Admin, error) {
	dtoList, err := a.adminRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	admins := make([]domain.Admin, len(dtoList))
	for i, dto := range dtoList {
		admins[i] = domain.Admin{
			ID:    dto.ID,
			UUID:  dto.UUID,
			Name:  dto.Name,
			Email: dto.Email,
		}
	}
	return admins, nil
}

func (a *AdminService) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	dto, err := a.adminRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &domain.Admin{
		ID:       dto.ID,
		UUID:     dto.UUID,
		Name:     dto.Name,
		Password: dto.Password,
		Email:    dto.Email,
	}, nil
}
