package domain

import "time"

type PostStatus int

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

type Media struct {
	Url  string
	Type string
}

type LinkHistory struct {
	ID               int       `json:"id"`
	UUID             string    `json:"uuid"`
	InstagramMediaID int       `json:"instagram_media_id"`
	InstagramLink    string    `json:"instagram_link"`
	WordpressMediaID int       `json:"wordpress_media_id"`
	WordpressLink    string    `json:"wordpress_link"`
	CreateAt         time.Time `json:"create_at"`
}
