package presenter

import (
	"errors"
	"net/http"
)

type Presenter struct {
}

func NewPresenter() Presenter {
	return Presenter{}
}

func (p *Presenter) Generate(err error, body any) (int, interface{}) {
	if err != nil {
		return http.StatusOK, body
	}
	switch {
	case errors.Is(err, ErrAuthentication):
		return http.StatusUnauthorized, err
	case errors.Is(err, ErrAuthorization):
		return http.StatusUnauthorized, err
	}
	return http.StatusNotImplemented, nil
}
