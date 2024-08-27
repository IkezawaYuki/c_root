package entity

import "time"

type Customer struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Password       string     `json:"password"`
	Email          string     `json:"email"`
	WordpressURL   string     `json:"wordpress_url"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag int        `json:"delete_hash_flag"`
}
