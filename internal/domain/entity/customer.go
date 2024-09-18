package entity

import "time"

type Customer struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Password       string     `json:"password"`
	Email          string     `json:"email"`
	WordpressURL   string     `json:"wordpressUrl"`
	FacebookToken  *string    `json:"facebookToken"`
	StartDate      *time.Time `json:"startDate"`
	InstagramID    *string    `json:"instagramId"`
	InstagramName  *string    `json:"instagramName"`
	DeleteHashFlag int        `json:"deleteHashFlag"`
}
