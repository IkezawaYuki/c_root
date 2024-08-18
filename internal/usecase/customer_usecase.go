package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/IkezawaYuki/c_root/internal/service"
)

type CustomerUsecase struct {
	baseRepository   repository.BaseRepository
	customerService  service.CustomerService
	authService      service.AuthService
	wordpressRestApi service.WordpressRestAPI
	graphApi         service.GraphAPI
	fileTransfer     service.FileTransfer
}

func NewCustomerUsecase(baseRepo repository.BaseRepository,
	customerSrv service.CustomerService,
	authSrv service.AuthService,
	wordpressRestApi service.WordpressRestAPI,
	graphApi service.GraphAPI,
) CustomerUsecase {
	return CustomerUsecase{
		baseRepository:   baseRepo,
		customerService:  customerSrv,
		authService:      authSrv,
		wordpressRestApi: wordpressRestApi,
		graphApi:         graphApi,
	}
}

func (c *CustomerUsecase) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return c.customerService.GetCustomer(ctx, id)
}

func (c *CustomerUsecase) Login(ctx context.Context, user *domain.User) (string, error) {
	customer, err := c.customerService.GetCustomerByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer); err != nil {
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
		detail, err := c.graphApi.GetMediaDetail(ctx, *customer.FacebookToken, mediaID)
		if err != nil {
			return err
		}

		localPathList, err := c.graphApi.DownloadMedias(ctx, *customer.FacebookToken, detail)
		if err != nil {
			return err
		}

		if err := c.customerService.SaveInstagramPost(ctx, detail, customer.StartDate); err != nil {
			return err
		}
		notYetPosts, err := c.customerService.GetInstagramPostNotYet(ctx, customer.ID)
		if err != nil {
			return err
		}
		for _, post := range notYetPosts {

		}
	}
	return err
}

func (c *CustomerUsecase) PostToWordpress(ctx context.Context, customerID string) error {
	customer, err := c.customerService.GetCustomer(ctx, customerID)
	if err != nil {
		return err
	}
	notYetPosts, err := c.customerService.GetInstagramPostNotYet(ctx, customer.ID)
	if err != nil {
		return err
	}
	for _, instagram := range notYetPosts {
		if err := c.fileTransfer.MakeTempDirectory(); err != nil {
			return err
		}

		//if detail.MediaType == "CAROUSEL_ALBUM" {
		//
		//	for _, child := range detail.Children {
		//
		//	}
		//} else if detail.MediaType == "VIDEO" {
		//
		//} else if detail.MediaType == "IMAGE" {
		//
		//}

		// インスタグラムの画像をローカルにダウンロードする
		//filePathList, err := c.graphApi.DownloadMedias(detail)
		//if err != nil {
		//	return err
		//}
		uploadFiles := make([]string, 0, len(filePathList))
		for _, filePath := range filePathList {
			// // ローカルのファイルをWordpressにファイルアップロードする
			c.wordpressRestApi.UploadFile()
		}

		// // ローカルのファイル削除する
		if err := c.fileTransfer.RemoveTempDirectory(); err != nil {
			return err
		}

		// Wordpressに投稿する
		c.wordpressRestApi.Post(ctx, instagram)

		// DBのレコードを投稿済みに更新する
		c.customerService.UpdateInstagramPost(ctx)

	}
	return nil
}
