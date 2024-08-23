package entity

import (
	"net/url"
	"strings"
	"time"
)

type InstagramPost struct {
	ID              string
	Caption         string
	MediaType       string
	MediaURL        string
	Timestamp       time.Time
	ChildrenID      []string
	ChildrenContent []ChildMedia
}

type ChildMedia struct {
	ID        string
	MediaURL  string
	MediaType string
}

func (i *InstagramPost) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}

func (i *InstagramPost) Title() string {
	return strings.Split(i.Caption, " ")[0]
}

func (i *ChildMedia) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}
