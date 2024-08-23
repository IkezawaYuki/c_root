package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
)

type CustomerUsecase struct {
	baseRepository     *repository.BaseRepository
	customerService    *service.CustomerService
	authService        *service.AuthService
	linkHistoryService *service.LinkHistoryService
	wordpressRestApi   *service.WordpressRestAPI
	graphApi           *service.GraphAPI
	fileTransfer       *service.FileService
}

func NewCustomerUsecase(
	baseRepo *repository.BaseRepository,
	customerSrv *service.CustomerService,
	authSrv *service.AuthService,
	linkHistoryService *service.LinkHistoryService,
	wordpressRestApi *service.WordpressRestAPI,
	graphApi *service.GraphAPI,
	fileTransfer *service.FileService,
) *CustomerUsecase {
	return &CustomerUsecase{
		baseRepository:     baseRepo,
		customerService:    customerSrv,
		authService:        authSrv,
		linkHistoryService: linkHistoryService,
		wordpressRestApi:   wordpressRestApi,
		graphApi:           graphApi,
		fileTransfer:       fileTransfer,
	}
}

func (c *CustomerUsecase) FindAll(ctx context.Context) ([]domain.Customer, error) {
	return c.customerService.FindAll(ctx)
}

func (c *CustomerUsecase) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return c.customerService.GetCustomer(ctx, id)
}

func (c *CustomerUsecase) Login(ctx context.Context, user *domain.User) (string, error) {
	customer, err := c.customerService.GetCustomerByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer.Password); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return c.authService.GenerateJWTCustomer(customer)
}

func (c *CustomerUsecase) FetchInstagramMediaFromGraphAPI(ctx context.Context, customerID string) error {
	customer, err := c.customerService.GetCustomer(ctx, customerID)
	if err != nil {
		return err
	}
	if customer.FacebookToken == nil {
		return errors.New("invalid operation")
	}
	mediaList, err := c.graphApi.GetMediaList(ctx, *customer.FacebookToken, *customer.InstagramID)
	if err != nil {
		return err
	}
	for _, mediaID := range mediaList {
		instaDetail, err := c.graphApi.GetMediaDetail(ctx, *customer.FacebookToken, mediaID)
		if err != nil {
			return err
		}
		medias, err := c.graphApi.FetchMedias(ctx, *customer.FacebookToken, instaDetail)
		if err != nil {
			return err
		}
		if err := c.customerService.SaveInstagramPost(ctx, instaDetail, medias, customer.StartDate); err != nil {
			return err
		}
	}
	return err
}

func (c *CustomerUsecase) PostToWordpress(ctx context.Context, customerID string) error {
	customer, err := c.customerService.GetCustomer(ctx, customerID)
	if err != nil {
		return err
	}
	notYetPosts, err := c.customerService.GetInstagramPostNotYet(ctx, customer.UUID)
	if err != nil {
		return err
	}
	if err := c.fileTransfer.MakeTempDirectory(customer.UUID); err != nil {
		return err
	}
	for _, instagram := range notYetPosts {
		instaDetail, err := c.graphApi.GetMediaDetail(ctx, *customer.FacebookToken, instagram.UUID)
		if err != nil {
			return err
		}
		medias, err := c.graphApi.FetchMedias(ctx, *customer.FacebookToken, instaDetail)
		if err != nil {
			return err
		}
		if err := c.customerService.SaveInstagramPost(ctx, instaDetail, medias, customer.StartDate); err != nil {
			return err
		}
		targetFiles, err := c.graphApi.FetchMedias(ctx, *customer.FacebookToken, instaDetail)
		if err != nil {
			return err
		}
		localPathList, err := c.fileTransfer.DownloadMedias(ctx, targetFiles)
		if err != nil {
			return err
		}
		wpMedia, err := c.wordpressRestApi.UploadFiles(ctx, customer.WordpressURL, localPathList)
		if err != nil {
			return err
		}
		wpLink, err := c.wordpressRestApi.CreatePosts(ctx, customer.WordpressURL, instaDetail, wpMedia)
		if err != nil {
			return err
		}
		if err := c.customerService.CreateInstagramWordpress(ctx, instaDetail.Permalink, wpLink); err != nil {
			return err
		}
	}
	if err := c.fileTransfer.RemoveTempDirectory(); err != nil {
		return err
	}
	return nil
}
