package domain

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
