package usecase

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
)

type CustomerUsecase struct {
	baseRepository   *repository.BaseRepository
	customerService  *service.CustomerService
	authService      *service.AuthService
	postService      *service.PostService
	wordpressRestApi *service.WordpressRestAPI
	graphApi         *service.GraphAPI
	fileTransfer     *service.FileService
}

func NewCustomerUsecase(
	baseRepo *repository.BaseRepository,
	customerSrv *service.CustomerService,
	authSrv *service.AuthService,
	postService *service.PostService,
	wordpressRestApi *service.WordpressRestAPI,
	graphApi *service.GraphAPI,
	fileTransfer *service.FileService,
) *CustomerUsecase {
	return &CustomerUsecase{
		baseRepository:   baseRepo,
		customerService:  customerSrv,
		authService:      authSrv,
		postService:      postService,
		wordpressRestApi: wordpressRestApi,
		graphApi:         graphApi,
		fileTransfer:     fileTransfer,
	}
}

func (c *CustomerUsecase) FindAll(ctx context.Context) ([]entity.Customer, error) {
	return c.customerService.FindAll(ctx)
}

func (c *CustomerUsecase) GetCustomer(ctx context.Context, id int) (*entity.Customer, error) {
	return c.customerService.FindByID(ctx, id)
}

func (c *CustomerUsecase) Login(ctx context.Context, user *entity.User) (string, error) {
	customer, err := c.customerService.GetCustomerByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer.Password); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return c.authService.GenerateJWTCustomer(customer)
}

func (c *CustomerUsecase) FetchAndPost(ctx context.Context, customerID int) error {
	customer, err := c.customerService.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if customer.FacebookToken == nil {
		return fmt.Errorf("customer.FacebookToken is nil")
	}
	mediaList, err := c.graphApi.GetMediaIDList(ctx, customer.FacebookToken, customer.InstagramID)
	if err != nil {
		return err
	}
	for _, media := range mediaList {
		linked, err := c.postService.IsLinked(ctx, media)
		if err != nil {
			return err
		}
		if linked {
			continue
		}
		detail, err := c.graphApi.GetMediaDetail(ctx, customer.FacebookToken, media)
		if err != nil {
			return err
		}
		post, err := c.postService.SaveInstagramPost(ctx, customerID, detail)
		if err != nil {
			return err
		}
		if err := c.graphApi.GetMediaChild(ctx, customer.FacebookToken, detail); err != nil {
			return err
		}
		if err := c.fileTransfer.MakeTempDirectory(customerID); err != nil {
			return err
		}
		mediaPaths, err := c.fileTransfer.DownloadMedias(ctx, detail)
		if err != nil {
			return err
		}
		wordpressMedia, err := c.wordpressRestApi.UploadFiles(ctx, customer.WordpressURL, mediaPaths)
		if err != nil {
			return err
		}
		wordpressLink, err := c.wordpressRestApi.CreatePosts(ctx, customer.WordpressURL, detail, wordpressMedia)
		if err != nil {
			return err
		}
		post.WordpressLink = &wordpressLink
		if err := c.postService.SaveWordpressPost(ctx, post); err != nil {
			return err
		}
	}
	return nil
}

func (c *CustomerUsecase) FetchAndPostByInstagramID(ctx context.Context, customerID int, instagramID string) error {
	customer, err := c.customerService.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if customer.FacebookToken == nil {
		return fmt.Errorf("customer.FacebookToken is nil")
	}
	detail, err := c.graphApi.GetMediaDetail(ctx, customer.FacebookToken, instagramID)
	if err != nil {
		return err
	}
	post, err := c.postService.SaveInstagramPost(ctx, customerID, detail)
	if err != nil {
		return err
	}
	if err := c.graphApi.GetMediaChild(ctx, customer.FacebookToken, detail); err != nil {
		return err
	}
	if err := c.fileTransfer.MakeTempDirectory(customerID); err != nil {
		return err
	}
	mediaPaths, err := c.fileTransfer.DownloadMedias(ctx, detail)
	if err != nil {
		return err
	}
	wordpressMedia, err := c.wordpressRestApi.UploadFiles(ctx, customer.WordpressURL, mediaPaths)
	if err != nil {
		return err
	}
	wordpressLink, err := c.wordpressRestApi.CreatePosts(ctx, customer.WordpressURL, detail, wordpressMedia)
	if err != nil {
		return err
	}
	post.WordpressLink = &wordpressLink
	if err := c.postService.SaveWordpressPost(ctx, post); err != nil {
		return err
	}
	return nil
}
