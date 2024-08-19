package service

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type FileService struct {
	httpClient *infrastructure.HttpClient
}

func NewFileService(httpClient *infrastructure.HttpClient) *FileService {
	return &FileService{
		httpClient: httpClient,
	}
}

func (f *FileService) DownloadMedias(ctx context.Context, medias []domain.Media) ([]string, error) {
	var result []string
	for _, media := range medias {
		path, err := f.DownloadMedia(ctx, media)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}
	return result, nil
}

func (f *FileService) DownloadMedia(ctx context.Context, media domain.Media) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", media.Url, nil)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(media.Url)
	filePath := filepath.Join(tempDirectory, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

const tempDirectory = "./tmp"

func (f *FileService) MakeTempDirectory(customerID string) error {
	err := os.Mkdir(tempDirectory+"_"+customerID, 0777)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (f *FileService) RemoveTempDirectory() error {
	return os.RemoveAll(tempDirectory)
}
