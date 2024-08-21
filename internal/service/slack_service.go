package service

import (
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type SlackService struct {
	webhookURL string
	httpClient *infrastructure.HttpClient
}

func NewSlackService(httpClient *infrastructure.HttpClient) *SlackService {
	return &SlackService{
		webhookURL: config.Env.SlackWebhookURL,
		httpClient: httpClient,
	}
}

func (s *SlackService) SendAlert(msg string) {

}

func (s *SlackService) SendNotification(msg string) {

}
