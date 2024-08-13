package presenter

import "net/http"

type HTTPError struct {
	StatusMessage string `json:"statusMessage"`
	ErrorCode     string `json:"errorCode"`
	Abstraction   string `json:"abstraction"`
}

type HTTPErrorPresenter interface {
	GenError(status int, errCode string, err error) (int, *HTTPError)
}

func NewHTTPErrorPresenter() HTTPErrorPresenter {
	return &httpErrorPresenter{}
}

type httpErrorPresenter struct{}

func (h httpErrorPresenter) GenError(status int, errCode string, err error) (int, *HTTPError) {
	return status, &HTTPError{
		StatusMessage: http.StatusText(status),
		ErrorCode:     errCode,
		Abstraction:   err.Error(),
	}
}
