package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type GraphAPI struct {
	httpClient *infrastructure.HttpClient
	baseURL    string
}

func NewGraph(httpClient *infrastructure.HttpClient) *GraphAPI {
	return &GraphAPI{
		httpClient: httpClient,
		baseURL:    config.Env.GraphApiURL,
	}
}

const getMediaChildURL = "/%s?fields=media_url,media_type"

func (i *GraphAPI) getMediaChild(ctx context.Context, facebookToken string, mediaID string) (*domain.InstagramMediaContent, error) {
	resp, err := i.httpClient.GetRequest(ctx, i.baseURL+fmt.Sprintf(getMediaChildURL, mediaID), fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return nil, err
	}
	var content domain.InstagramMediaContent
	if err := json.Unmarshal(resp, &content); err != nil {
		return nil, err
	}
	return &content, nil
}

func (i *GraphAPI) FetchMedias(ctx context.Context, facebookToken string, detail *domain.InstagramMediaDetail) ([]domain.Media, error) {
	if detail.MediaType == "CAROUSEL_ALBUM" {
		medias := make([]domain.Media, len(detail.Children))
		for _, child := range detail.Children {
			cMedia, err := i.getMediaChild(ctx, facebookToken, child)
			if err != nil {
				return nil, err
			}
			medias = append(medias, domain.Media{
				ID:   cMedia.ID,
				Url:  cMedia.MediaURL,
				Type: cMedia.MediaType,
			})
		}
		return medias, nil
	}
	medias := make([]domain.Media, 1)
	medias[0] = domain.Media{
		ID:   detail.ID,
		Url:  detail.MediaURL,
		Type: detail.MediaType,
	}
	return medias, nil
}

const getInstagramBusinessAccountURL = "/me?fields=id,name,accounts{instagram_business_account}"

func (i *GraphAPI) GetInstagramBusinessAccountID(ctx context.Context, facebookToken string) (string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+getInstagramBusinessAccountURL,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return "", err
	}
	var instagram domain.GraphApiMeResponse
	err = json.Unmarshal(resp, &instagram)
	if err != nil {
		return "", err
	}
	return instagram.InstagramBusinessAccountID(), nil
}

const getMediaList = "/%s"

func (i *GraphAPI) GetMediaList(ctx context.Context, facebookToken, instagramID string) ([]string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+fmt.Sprintf(getMediaList, instagramID),
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaList
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return detail.ConvertToInstagramMediaList(), nil
}

const getMediaContent = "/%s?fields=media_type,media_url,id,caption,timestamp,permalink,children"

func (i *GraphAPI) GetMediaDetail(ctx context.Context, facebookToken string, mediaID string) (*domain.InstagramMediaDetail, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+fmt.Sprintf(getMediaContent, mediaID),
		fmt.Sprintf("Bearer %s", facebookToken),
	)
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaDetail
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}
