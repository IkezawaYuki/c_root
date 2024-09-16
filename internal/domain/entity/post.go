package entity

import "time"

type Post struct {
	ID               int       `json:"id"`
	CustomerID       int       `json:"customer_id"`
	InstagramMediaID string    `json:"instagram_media_id"`
	InstagramLink    string    `json:"instagram_link"`
	WordpressLink    *string   `json:"wordpress_link"`
	WordpressMediaID *string   `json:"wordpress_media_id"`
	CreatedAt        time.Time `json:"created_at"`
}
