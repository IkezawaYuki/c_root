package interfaces

import (
	"context"
	"github.com/IkezawaYuki/c_root/internal/croot/domain"
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre/driver"
	"io"
)

type instagramAdapter struct {
	client driver.HttpClient
}

func NewInstagramAdapter(client driver.HttpClient) domain.InstagramPostRepository {
	return &instagramAdapter{client: client}
}

const endpoint = "https://graph.facebook.com/v18.0"

func (i *instagramAdapter) FindAll(ctx context.Context, accessToken string) (*domain.InstagramMediaList, error) {

}

func (i *instagramAdapter) FindDetail(ctx context.Context, accessToken string, id string) (*domain.InstagramMediaDetail, error) {

}

func (i *instagramAdapter) DownloadMedia(ctx context.Context, accessToken string, id string) (io.ReadCloser, error) {

}
