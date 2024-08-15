package di

import (
	"github.com/IkezawaYuki/c_root/internal/controller"
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/IkezawaYuki/c_root/internal/service"
	"github.com/IkezawaYuki/c_root/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAuthService(db *gorm.DB, redisCli *redis.Client) service.AuthService {
	customerRepo := repository.NewCustomerRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	return service.NewAuthService(customerRepo, redisClient)
}

func NewCustomerService(db *gorm.DB, redisCli *redis.Client) service.CustomerService {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	return service.NewCustomerService(baseRepo, customerRepo)
}

func NewCustomerController(db *gorm.DB, redisCli *redis.Client) controller.CustomerController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(baseRepo, customerRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	customerUsecase := usecase.NewCustomerUsecase(baseRepo, customerService, authService)
	return controller.NewCustomerController(customerUsecase, pre)
}

func NewAdminController(db *gorm.DB, redisCli *redis.Client) controller.AdminController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(baseRepo, customerRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	adminService := service.NewAdminService(customerRepo, adminRepo)
	adminUsecase := usecase.NewAdminUsecase(baseRepo, adminService, authService, customerService)
	return controller.NewAdminController(adminUsecase, pre)
}
