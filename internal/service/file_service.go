package service

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileService struct {
	httpClient *infrastructure.HttpClient
}

func NewFileService(httpClient *infrastructure.HttpClient) *FileService {
	return &FileService{
		httpClient: httpClient,
	}
}

func (f *FileService) DownloadMedias(ctx context.Context, customerID int, post *entity.InstagramPost) ([]string, error) {
	var fileList []string
	if len(post.ChildrenContent) == 0 {
		mediaPath, err := f.DownloadMedia(ctx, customerID, post.MediaURL)
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, mediaPath)
	} else {
		for _, child := range post.ChildrenContent {
			mediaPath, err := f.DownloadMedia(ctx, customerID, child.MediaURL)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, mediaPath)
		}
	}
	return fileList, nil
}

func (f *FileService) DownloadMedia(ctx context.Context, customerID int, mediaUrl string) (string, error) {
	fmt.Println("DownloadMedia is invoked")
	fmt.Println(mediaUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", mediaUrl, nil)
	if err != nil {
		return "", err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	filename := strings.Split(filepath.Base(mediaUrl), "?")[0]
	filePath := filepath.Join(fmt.Sprintf(tempDirectory, customerID), filename)
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

const tempDirectory = "./tmp_%d"

func (f *FileService) MakeTempDirectory(customerID int) error {
	err := os.Mkdir(fmt.Sprintf(tempDirectory, customerID), 0777)
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
