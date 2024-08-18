package domain

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"net/url"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("email already in use")
	ErrNotFound       = errors.New("not found")
	ErrAuthentication = errors.New("authentication err")
	ErrAuthorization  = errors.New("authorization err")
	ErrDuplicateKey   = errors.New("duplicate key err")
)

type Customer struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Password       string     `json:"password"`
	Email          string     `json:"email"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag int        `json:"delete_hash_flag"`
}

type CustomerDto struct {
	ID             string
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	FacebookToken  sql.NullString
	StartDate      sql.NullTime
	InstagramID    sql.NullString
	InstagramName  sql.NullString
	DeleteHashFlag int
	gorm.Model
}

func (c *CustomerDto) TableName() string {
	return "customers"
}

func (c *CustomerDto) ConvertToCustomer() *Customer {
	var customer Customer
	customer.ID = c.ID
	customer.Name = c.Name
	customer.Email = c.Email
	customer.Password = c.Password
	if c.FacebookToken.Valid {
		customer.FacebookToken = &c.FacebookToken.String
	}
	if c.StartDate.Valid {
		customer.StartDate = &c.StartDate.Time
	}
	if c.InstagramID.Valid {
		customer.InstagramID = &c.InstagramID.String
	}
	if c.InstagramName.Valid {
		customer.InstagramName = &c.InstagramName.String
	}
	return &customer
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	return getFilename(i.MediaURL)
}

type InstagramMediaContent struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
}

func (i *InstagramMediaContent) FileName() (string, error) {
	return getFilename(i.MediaURL)
}

func getFilename(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}

type PostStatus int

type InstagramDto struct {
	ID         string
	Caption    string
	MediaType  string
	MediaURL   string
	Permalink  string
	PostStatus int
	Timestamp  time.Time
	gorm.Model
}

func (i *InstagramDto) TableName() string { return "instagram_posts" }

type InstagramPost struct {
	ID         string
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

type GraphApiMeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Accounts struct {
		Data []struct {
			InstagramBusinessAccount struct {
				ID string `json:"id"`
			} `json:"instagram_business_account"`
			ID string `json:"id"`
		} `json:"data"`
		Paging struct {
			Cursors struct {
				Before string `json:"before"`
				After  string `json:"after"`
			} `json:"cursors"`
		} `json:"paging"`
	} `json:"accounts"`
}

func (r *GraphApiMeResponse) InstagramBusinessAccountID() string {
	return r.Accounts.Data[0].InstagramBusinessAccount.ID
}
