package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (a *AuthService) IsCustomerIsLogin(tokenString string) (int, error) {
	slog.Info("IsCustomerIsLogin is invoked")
	slog.Info(tokenString)
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		slog.Info(err.Error())
		return 0, objects.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, objects.ErrAuthorization
	}
	if !claims.VerifyAudience("customer", true) {
		return 0, objects.ErrAuthentication
	}

	return int(claims["sub"].(float64)), nil
}

func (a *AuthService) IsAdminLogin(tokenString string) (int, error) {
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
		return 0, objects.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, objects.ErrAuthorization
	}
	if !claims.VerifyAudience("admin", true) {
		return 0, objects.ErrAuthentication
	}
	return int(claims["sub"].(float64)), nil
}

func (a *AuthService) GenerateJWTCustomer(c *entity.Customer) (string, error) {
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

func (a *AuthService) GenerateJWTAdmin(admin *entity.Admin) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "admin",
		"sub":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) CheckPassword(user *entity.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return fmt.Errorf("password is incorrect: %s, %v", err.Error(), objects.ErrAuthorization)
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

func (a *AdminService) GetCustomerByID(ctx context.Context, id int) (*entity.Customer, error) {
	customerModel, err := a.customerRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var customer entity.Customer
	customer.ID = int(customerModel.ID)
	customer.Name = customerModel.Name
	customer.Password = customerModel.Password
	customer.Email = customerModel.Email
	customer.WordpressURL = customerModel.WordpressURL
	if customerModel.FacebookToken.Valid {
		customer.FacebookToken = &customerModel.FacebookToken.String
	}
	if customerModel.StartDate.Valid {
		customer.StartDate = &customerModel.StartDate.Time
	}
	if customerModel.InstagramID.Valid {
		customer.InstagramID = &customerModel.InstagramID.String
	}
	if customerModel.InstagramName.Valid {
		customer.InstagramName = &customerModel.InstagramName.String
	}
	customer.DeleteHashFlag = customerModel.DeleteHashFlag
	return &customer, nil
}

func (a *AdminService) FindAll(ctx context.Context) ([]entity.Admin, error) {
	modelList, err := a.adminRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	admins := make([]entity.Admin, len(modelList))
	for i, m := range modelList {
		admins[i] = entity.Admin{
			ID:       m.ID,
			Name:     m.Name,
			Password: m.Password,
			Email:    m.Email,
		}
	}
	return admins, nil
}

func (a *AdminService) FindByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	m, err := a.adminRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &entity.Admin{
		ID:       m.ID,
		Name:     m.Name,
		Password: m.Password,
		Email:    m.Email,
	}, nil
}

func (a *AdminService) FindByID(ctx context.Context, id int) (*entity.Admin, error) {
	m, err := a.adminRepository.FindById(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	return &entity.Admin{
		ID:       m.ID,
		Name:     m.Name,
		Password: m.Password,
		Email:    m.Email,
	}, nil
}

func (a *AdminService) CreateAdmin(ctx context.Context, admin *entity.Admin) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	adminModel := model.Admin{
		Name:     admin.Name,
		Email:    admin.Email,
		Password: string(passwordHash),
	}
	if err := a.adminRepository.Save(ctx, &adminModel); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return objects.ErrDuplicateEmail
		}
		return err
	}
	admin.ID = adminModel.ID
	return nil
}
