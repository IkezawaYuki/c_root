package domain

import (
	"time"
)

type Customer struct {
	ID             uint       `json:"id"`
	UUID           string     `json:"uuid"`
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
