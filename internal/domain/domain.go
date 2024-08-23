package domain

import "time"

type PostStatus int

type Media struct {
	ID   string
	Url  string
	Type string
}

type LinkHistory struct {
	ID               int       `json:"id"`
	CustomerID       string    `json:"customer_id"`
	InstagramMediaID string    `json:"instagram_media_id"`
	InstagramLink    string    `json:"instagram_link"`
	WordpressMediaID string    `json:"wordpress_media_id"`
	WordpressLink    string    `json:"wordpress_link"`
	CreateAt         time.Time `json:"create_at"`
}
