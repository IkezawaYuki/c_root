package di

import (
	"github.com/IkezawaYuki/popple/internal/controller"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAuthService(db *gorm.DB, redisCli *redis.Client) *service.AuthService {
	customerRepo := repository.NewCustomerRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	return service.NewAuthService(customerRepo, redisClient)
}

func NewCustomerService(db *gorm.DB) *service.CustomerService {
	customerRepo := repository.NewCustomerRepository(db)
	postRepo := repository.NewPostRepository(db)
	return service.NewCustomerService(customerRepo, postRepo)
}

func NewCustomerController(db *gorm.DB, redisCli *redis.Client) controller.CustomerController {
	pre := presenter.NewPresenter()
	customerUsecase := NewCustomerUsecase(db, redisCli)
	return controller.NewCustomerController(customerUsecase, pre)
}

func NewAdminController(db *gorm.DB, redisCli *redis.Client) controller.AdminController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	postRepo := repository.NewPostRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(customerRepo, postRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	adminService := service.NewAdminService(customerRepo, adminRepo)
	customerUsecase := NewCustomerUsecase(db, redisCli)
	adminUsecase := usecase.NewAdminUsecase(baseRepo, adminService, authService, customerService, customerUsecase)
	return controller.NewAdminController(adminUsecase, pre)
}

func NewBatchController(db *gorm.DB, redisCli *redis.Client) controller.BatchController {
	pre := presenter.NewPresenter()
	httpClient := infrastructure.NewHttpClient()
	slack := service.NewSlackService(httpClient)
	customerUsecase := NewCustomerUsecase(db, redisCli)
	batchUsecase := usecase.NewBatchUsecase(customerUsecase, slack)
	return controller.NewBatchController(batchUsecase, pre)
}

func NewCustomerUsecase(db *gorm.DB, redisCli *redis.Client) *usecase.CustomerUsecase {
	httpClient := infrastructure.NewHttpClient()
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	postRepo := repository.NewPostRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	customerService := service.NewCustomerService(customerRepo, postRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	postService := service.NewPostService(postRepo)
	wordpressRestApi := service.NewWordpressRestAPI(httpClient)
	graphApi := service.NewGraph(httpClient)
	fileTransfer := service.NewFileService(httpClient)
	return usecase.NewCustomerUsecase(
		baseRepo,
		customerService,
		authService,
		postService,
		wordpressRestApi,
		graphApi,
		fileTransfer)
}
