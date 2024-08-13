package domain

import "context"

type InstagramMediaList struct {
	Media struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	} `json:"media"`
	ID string `json:"id"`
}

type InstagramMediaDetail struct {
	Caption   string `json:"caption"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
	Timestamp string `json:"timestamp"`
	ID        string `json:"id"`
}

type InstagramPostRepository interface {
	FindAll(ctx context.Context, accessToken string) (*InstagramMediaList, error)
	FindDetail(ctx context.Context, accessToken string, id string) (*InstagramMediaDetail, error)
}
