package entity

import "time"

type Post struct {
	ID               int       `json:"id"`
	CustomerID       int       `json:"customer_id"`
	InstagramMediaID string    `json:"instagram_media_id"`
	InstagramLink    *string   `json:"instagram_link"`
	WordpressLink    *string   `json:"wordpress_link"`
	CreatedAt        time.Time `json:"created_at"`
}
