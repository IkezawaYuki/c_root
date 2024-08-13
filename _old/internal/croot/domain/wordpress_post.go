package domain

type WordpressPostRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        string `json:"status"`
	FeaturedMedia string `json:"featured_media"`
}

type WordPressPostResponse struct {
	Link string `json:"link"`
}

type WordpressRepository interface {
	Post(request WordpressPostRequest) WordPressPostResponse
	UploadFile()
}
