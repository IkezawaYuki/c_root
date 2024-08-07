package entity

import "time"

type Customer struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag int        `json:"delete_hash_flag"`
}
