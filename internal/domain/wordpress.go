package domain

import (
	"fmt"
	"strings"
)

type WordpressPosts struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        string `json:"status"`
	FeaturedMedia string `json:"featured_media"`
}

type WordpressMedia struct {
	ID        string `json:"id"`
	SourceURL string `json:"source_url"`
	MediaType string `json:"media_type"`
}

func NewWordpressPosts(instaDetail *InstagramMediaDetail, wpMedia []*WordpressMedia) WordpressPosts {
	wordpressPosts := WordpressPosts{}
	wordpressPosts.Title = instaDetail.Title()
	wordpressPosts.FeaturedMedia = wpMedia[0].ID
	if instaDetail.MediaType == "IMAGE" {
		wordpressPosts.Content = fmt.Sprintf("%s%s", getImageHtml(wpMedia[0].SourceURL), getContentHtml(instaDetail.Caption))
	} else if instaDetail.MediaType == "VIDEO" {
		wordpressPosts.Content = fmt.Sprintf("%s%s", getVideoHtml(wpMedia[0].SourceURL), getContentHtml(instaDetail.Caption))
	} else {
		wordpressPosts.Content = getCarousel(instaDetail, wpMedia)
	}
	wordpressPosts.Status = "publish"
	return wordpressPosts
}

func getCarousel(instaDetail *InstagramMediaDetail, wpMedia []*WordpressMedia) string {
	sb := strings.Builder{}
	sb.WriteString("<div class='a-root-wordpress-instagram-slider'>")
	for _, media := range wpMedia {
		if media.MediaType == "IMAGE" {
			sb.WriteString(getImageHtml(media.SourceURL))
		} else if media.MediaType == "VIDEO" {
			sb.WriteString(getVideoHtml(media.SourceURL))
		}
	}
	sb.WriteString("</div>")
	sb.WriteString(getContentHtml(instaDetail.Caption))
	return sb.String()
}

func getVideoHtml(url string) string {
	return fmt.Sprintf(`<div><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>`, url)
}

func getImageHtml(url string) string {
	return fmt.Sprintf(`<div><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>
Sorry, your browser does not support embedded videos.</video></div>`, url)
}

func getContentHtml(caption string) string {
	sb := strings.Builder{}
	sb.WriteString("<p>")
	for i, row := range strings.Split(caption, "/n") {
		if i != 0 {
			sb.WriteString("<br>")
		}
		sb.WriteString(row)
	}
	sb.WriteString("</p>")
	return sb.String()
}
