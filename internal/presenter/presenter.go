package presenter

import (
	"errors"
	"github.com/IkezawaYuki/c_root/internal/domain"
	"log/slog"
	"net/http"
)

type Presenter struct {
}

func NewPresenter() Presenter {
	return Presenter{}
}

func (p *Presenter) Generate(err error, body any) (int, any) {
	slog.Info("Generate is invoked")
	if err == nil {
		return http.StatusOK, body
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrAuthentication):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, domain.ErrAuthorization):
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
