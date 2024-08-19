package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/infrastructure"
	"github.com/IkezawaYuki/c_root/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CustomerService struct {
	customerRepository  repository.CustomerRepository
	instagramRepository repository.InstagramRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository,
	instagramRepository repository.InstagramRepository,
) CustomerService {
	return CustomerService{
		customerRepository:  customerRepo,
		instagramRepository: instagramRepository,
	}
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.customerRepository.FindByID(ctx, id)
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	return s.customerRepository.FindByEmail(ctx, email)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	customer.ID = uuid.New().String()
	if err := s.customerRepository.Save(ctx, &domain.CustomerDto{
		ID:             customer.ID,
		Name:           customer.Name,
		Email:          customer.Email,
		Password:       string(passwordHash),
		DeleteHashFlag: 0,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *CustomerService) GetInstagramPostNotYet(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
	records, err := s.instagramRepository.FindNotYetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	posts := make([]domain.InstagramPost, len(records))
	for i, record := range records {
		posts[i] = domain.InstagramPost{
			ID:         record.ID,
			Caption:    record.Caption,
			MediaType:  record.Caption,
			MediaURL:   record.MediaURL,
			PostStatus: domain.PostStatus(record.PostStatus),
			Timestamp:  record.Timestamp,
		}
	}
	return posts, nil
}

func (s *CustomerService) GetInstagramPost(ctx context.Context, customerID string) ([]domain.InstagramPost, error) {
	records, err := s.instagramRepository.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	posts := make([]domain.InstagramPost, len(records))
	for i, record := range records {
		posts[i] = domain.InstagramPost{
			ID:         record.ID,
			Caption:    record.Caption,
			MediaType:  record.Caption,
			MediaURL:   record.MediaURL,
			PostStatus: domain.PostStatus(record.PostStatus),
			Timestamp:  record.Timestamp,
		}
	}
	return posts, nil
}

func (s *CustomerService) SaveInstagramPost(ctx context.Context, instagramPost *domain.InstagramMediaDetail, startDate *time.Time) error {
	timestamp, err := time.Parse("2006-01-02T15:04:05-0700", instagramPost.Timestamp)
	if err != nil {
		return err
	}
	if startDate == nil {
		return errors.New("startDate is required")
	}
	status := domain.NotYet
	if startDate.Before(timestamp) {
		status = domain.Linked
	}
	err = s.instagramRepository.Save(ctx, domain.InstagramDto{
		ID:         instagramPost.ID,
		Caption:    instagramPost.Caption,
		MediaType:  instagramPost.MediaType,
		MediaURL:   instagramPost.MediaURL,
		Permalink:  instagramPost.Permalink,
		PostStatus: int(status),
		Timestamp:  timestamp,
	})
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return nil
		}
		return err
	}
	return nil
}

type AuthService struct {
	customerRepository repository.CustomerRepository
	redisClient        repository.RedisClient
}

func NewAuthService(customerRepo repository.CustomerRepository, redisClient repository.RedisClient) AuthService {
	return AuthService{
		customerRepository: customerRepo,
		redisClient:        redisClient,
	}
}

func (a *AuthService) IsCustomerIsLogin(tokenString string) (string, error) {
	slog.Info("IsCustomerIsLogin is invoked")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		slog.Info(err.Error())
		return "", domain.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", domain.ErrAuthorization
	}
	if !claims.VerifyAudience("customer", true) {
		return "", domain.ErrAuthentication
	}
	return claims["sub"].(string), nil
}

func (a *AuthService) IsAdminLogin() (bool, error) {
	return true, nil
}

func (a *AuthService) GenerateJWTCustomer(c *domain.Customer) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "customer",
		"sub":   c.ID,
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *AuthService) CheckPassword(user *domain.User, customer *domain.Customer) error {
	return bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(user.Password))
}

type AdminService struct {
	customerRepository repository.CustomerRepository
	adminRepository    repository.AdminRepository
}

func NewAdminService(customerRepo repository.CustomerRepository, adminRepo repository.AdminRepository) AdminService {
	return AdminService{
		customerRepository: customerRepo,
		adminRepository:    adminRepo,
	}
}

func (a *AdminService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return a.customerRepository.FindByID(ctx, id)
}

type GraphAPI struct {
	httpClient infrastructure.HttpClient
}

func NewGraph(instagramRepo repository.InstagramRepository, httpClient infrastructure.HttpClient) GraphAPI {
	return GraphAPI{
		httpClient: httpClient,
	}
}

const getMediaChildURL = "%s?fields=media_url,media_type"

func (i *GraphAPI) getMediaChild(ctx context.Context, facebookToken string, mediaID string) (*domain.InstagramMediaContent, error) {
	resp, err := i.httpClient.GetRequest(ctx, fmt.Sprintf(getMediaChildURL, mediaID), nil, fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return nil, err
	}
	var content domain.InstagramMediaContent
	if err := json.Unmarshal(resp, &content); err != nil {
		return nil, err
	}
	return &content, nil
}

func (i *GraphAPI) FetchMedias(ctx context.Context, facebookToken string, detail *domain.InstagramMediaDetail) ([]domain.Media, error) {
	if detail.MediaType == "CAROUSEL_ALBUM" {
		medias := make([]domain.Media, len(detail.Children))
		for _, child := range detail.Children {
			cMedia, err := i.getMediaChild(ctx, facebookToken, child)
			if err != nil {
				return nil, err
			}
			medias = append(medias, domain.Media{
				Url:  cMedia.MediaURL,
				Type: cMedia.MediaType,
			})
		}
		return medias, nil
	}

	return []domain.Media{
		{detail.MediaType, detail.MediaURL},
	}, nil
}

const getInstagramBusinessAccountURL = "/me?fields=id,name,accounts{instagram_business_account}"

func (i *GraphAPI) GetInstagramBusinessAccountID(ctx context.Context, facebookToken string) (string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		getInstagramBusinessAccountURL,
		nil,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return "", err
	}
	var instagram domain.GraphApiMeResponse
	err = json.Unmarshal(resp, &instagram)
	if err != nil {
		return "", err
	}
	return instagram.InstagramBusinessAccountID(), nil
}

const getMediaList = "/%s"

func (i *GraphAPI) GetMediaList(ctx context.Context, facebookToken, instagramID string) ([]string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		fmt.Sprintf(getMediaList, instagramID),
		nil,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaList
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return detail.ConvertToInstagramMediaList(), nil
}

const getMediaContent = "/%s?fields=media_type,media_url,id,caption,timestamp,permalink,children"

func (i *GraphAPI) GetMediaDetail(ctx context.Context, facebookToken string, mediaID string) (*domain.InstagramMediaDetail, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		fmt.Sprintf(getMediaContent, mediaID),
		nil,
		fmt.Sprintf("Bearer %s", facebookToken),
	)
	if err != nil {
		return nil, err
	}
	var detail domain.InstagramMediaDetail
	if err := json.Unmarshal(resp, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

type WordpressRestAPI struct {
	httpClient infrastructure.HttpClient
}

func (w *WordpressRestAPI) CreatePosts(ctx context.Context, instaDetail *domain.InstagramMediaDetail, idList []string) error {
	posts := domain.NewWordpressPosts(instaDetail, idList)

}

func (w *WordpressRestAPI) UploadFiles(ctx context.Context, pathList []string) ([]string, error) {
	var mediaIDList []string
	for _, path := range pathList {
		mediaID, err := w.UploadFile(ctx, path)
		if err != nil {
			return nil, err
		}
		mediaIDList = append(mediaIDList, mediaID)
	}
	return mediaIDList, nil
}

const wpMediaUrl = "/wp-json/wp/v2/media"

type RespID struct {
	ID string `json:"id"`
}

func (w *WordpressRestAPI) UploadFile(ctx context.Context, path string) (string, error) {
	var respId RespID
	resp, err := w.httpClient.UploadFile(ctx, wpMediaUrl, path, "auth")
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(resp, &respId); err != nil {
		return "", err
	}
	return respId.ID, nil
}

type FileTransfer struct {
	httpClient infrastructure.HttpClient
}

func (f *FileTransfer) DownloadMedias(ctx context.Context, medias []domain.Media) ([]string, error) {
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

func (f *FileTransfer) DownloadMedia(ctx context.Context, media domain.Media) (string, error) {
	resp, err := http.Get(media.Url)
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

func (f *FileTransfer) MakeTempDirectory() error {
	err := os.Mkdir(tempDirectory, 0777)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (f *FileTransfer) RemoveTempDirectory() error {
	return os.RemoveAll(tempDirectory)
}
