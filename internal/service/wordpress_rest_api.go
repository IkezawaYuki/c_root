package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type WordpressRestAPI struct {
	httpClient *infrastructure.HttpClient
}

func NewWordpressRestAPI(httpClient *infrastructure.HttpClient) *WordpressRestAPI {
	return &WordpressRestAPI{
		httpClient: httpClient,
	}
}

const wpPostUrl = "/wp/v2/posts"

type CreatePostResponse struct {
	Link string `json:"link"`
}

func (w *WordpressRestAPI) CreatePosts(ctx context.Context, wordpressURL string, instaDetail *domain.InstagramMediaDetail, wpMedia []*domain.WordpressMedia) (string, error) {
	posts := domain.NewWordpressPosts(instaDetail, wpMedia)
	resp, err := w.httpClient.PostRequest(ctx, wordpressURL+wpPostUrl, posts, BasicAuthHeader())
	if err != nil {
		return "", err
	}
	var response CreatePostResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return "", err
	}
	return response.Link, nil
}

func (w *WordpressRestAPI) UploadFiles(ctx context.Context, wordpressURL string, pathList []string) ([]*domain.WordpressMedia, error) {
	var wordpressMediaList []*domain.WordpressMedia
	for _, path := range pathList {
		wordpressMedia, err := w.UploadFile(ctx, wordpressURL, path)
		if err != nil {
			return nil, err
		}
		wordpressMediaList = append(wordpressMediaList, wordpressMedia)
	}
	return wordpressMediaList, nil
}

const wpMediaUrl = "/wp-json/wp/v2/media"

func (w *WordpressRestAPI) UploadFile(ctx context.Context, wordpressURL string, path string) (*domain.WordpressMedia, error) {
	resp, err := w.httpClient.UploadFile(ctx, wordpressURL+wpMediaUrl, path, BasicAuthHeader())
	if err != nil {
		return nil, err
	}
	var wordpressMedia domain.WordpressMedia
	if err := json.Unmarshal(resp, &wordpressMedia); err != nil {
		return nil, err
	}
	return &wordpressMedia, nil
}

func BasicAuthHeader() string {
	auth := config.Env.WordpressAdminEmail + ":" + config.Env.WordpressAdminPassword
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encoded
}
