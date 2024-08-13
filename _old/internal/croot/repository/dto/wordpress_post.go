package dto

type WordPressPostRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        string `json:"status"`
	FeaturedMedia string `json:"featured_media"`
}

type WordPressPostResponse struct {
	Link string `json:"link"`
}
