package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"time"
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

const getInstagramBusinessAccountURL = "/me?fields=id,name,accounts{instagram_business_account}"

type GraphApiMeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Accounts struct {
		Data []struct {
			InstagramBusinessAccount struct {
				ID string `json:"id"`
			} `json:"instagram_business_account"`
			ID string `json:"id"`
		} `json:"data"`
	} `json:"accounts"`
}

func (i *GraphAPI) GetInstagramBusinessAccountID(ctx context.Context, facebookToken string) (string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+getInstagramBusinessAccountURL,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return "", err
	}
	var instagram GraphApiMeResponse
	err = json.Unmarshal(resp, &instagram)
	if err != nil {
		return "", err
	}
	if len(instagram.Accounts.Data) == 0 {
		return "", objects.ErrNotFound
	}
	return instagram.Accounts.Data[0].InstagramBusinessAccount.ID, nil
}

const getMediaList = "/%s"

type InstagramMediaList struct {
	Media struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	} `json:"media"`
	ID string `json:"id"`
}

func (i *GraphAPI) GetMediaIDList(ctx context.Context, facebookToken, instagramID *string) ([]string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+fmt.Sprintf(getMediaList, *instagramID),
		fmt.Sprintf("Bearer %s", *facebookToken))
	if err != nil {
		return nil, err
	}
	var mediaList InstagramMediaList
	if err := json.Unmarshal(resp, &mediaList); err != nil {
		return nil, err
	}
	mediaIdList := make([]string, len(mediaList.Media.Data))
	for idx, media := range mediaList.Media.Data {
		mediaIdList[idx] = media.ID
	}
	return mediaIdList, nil
}

const getMediaDetail = "/%s?fields=media_type,media_url,id,caption,timestamp,permalink,children"

type InstagramMediaDetail struct {
	ID        string   `json:"id"`
	Caption   string   `json:"caption"`
	MediaType string   `json:"media_type"`
	MediaURL  string   `json:"media_url"`
	Timestamp string   `json:"timestamp"`
	Permalink string   `json:"permalink"`
	Children  []string `json:"children"`
}

func (i *GraphAPI) GetMediaDetail(ctx context.Context, facebookToken *string, mediaID string) (*entity.InstagramPost, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+fmt.Sprintf(getMediaDetail, mediaID),
		fmt.Sprintf("Bearer %s", facebookToken),
	)
	if err != nil {
		return nil, err
	}
	var detail InstagramMediaDetail
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	var post entity.InstagramPost
	post.ID = detail.ID
	post.Caption = detail.Caption
	post.MediaType = detail.MediaType
	post.MediaURL = detail.MediaURL
	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", detail.Timestamp)
	if err != nil {
		return nil, err
	}
	post.Timestamp = timestamp

	if len(detail.Children) > 0 {
		children := make([]string, len(detail.Children))
		copy(children, detail.Children)
		post.ChildrenID = children
	}

	return &post, nil
}

const getMediaChildURL = "/%s?fields=media_url,media_type"

type InstagramMediaChild struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
}

func (i *GraphAPI) GetMediaChild(ctx context.Context, facebookToken *string, post *entity.InstagramPost) error {
	if len(post.ChildrenID) == 0 {
		return nil
	}
	contents := make([]entity.ChildMedia, len(post.ChildrenID))
	for idx, childID := range post.ChildrenID {
		resp, err := i.httpClient.GetRequest(ctx, i.baseURL+fmt.Sprintf(getMediaChildURL, childID), fmt.Sprintf("Bearer %s", facebookToken))
		if err != nil {
			return err
		}
		var content InstagramMediaChild
		if err := json.Unmarshal(resp, &content); err != nil {
			return err
		}
		contents[idx] = entity.ChildMedia{
			ID:        content.ID,
			MediaURL:  content.MediaURL,
			MediaType: content.MediaType,
		}
	}
	post.ChildrenContent = contents
	return nil
}
