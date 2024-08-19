package domain

import (
	"net/url"
	"strings"
	"time"
)

type InstagramPost struct {
	ID         uint
	UUID       string
	MediaID    int
	Caption    string
	MediaType  string
	MediaURL   string
	PostStatus PostStatus
	Timestamp  time.Time
}

var (
	NotYet PostStatus = 0
	Linked PostStatus = 1
)

type InstagramMediaList struct {
	Media struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	} `json:"media"`
	ID string `json:"id"`
}

func (i *InstagramMediaList) ConvertToInstagramMediaList() []string {
	var idList []string
	for _, media := range i.Media.Data {
		idList = append(idList, media.ID)
	}
	return idList
}

type InstagramMediaDetail struct {
	ID        string   `json:"id"`
	Caption   string   `json:"caption"`
	MediaType string   `json:"media_type"`
	MediaURL  string   `json:"media_url"`
	Timestamp string   `json:"timestamp"`
	Permalink string   `json:"permalink"`
	Children  []string `json:"children"`
}

func (i *InstagramMediaDetail) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}

func (i *InstagramMediaDetail) Title() string {
	return strings.Split(i.Caption, " ")[0]
}

type InstagramMediaContent struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
}

func (i *InstagramMediaContent) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}
