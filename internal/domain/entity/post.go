package entity

import "time"

type Post struct {
	ID               int       `json:"id"`
	CustomerID       int       `json:"customerId"`
	InstagramMediaID string    `json:"instagramMediaId"`
	InstagramLink    string    `json:"instagramLink"`
	WordpressLink    *string   `json:"wordpressLink"`
	WordpressMediaID *string   `json:"wordpressMediaId"`
	CreatedAt        time.Time `json:"createdAt"`
}
