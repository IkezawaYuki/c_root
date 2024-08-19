package presenter

import (
	"fmt"
	"github.com/IkezawaYuki/popple/internal/croot/domain/crooterrors"
	"net/http"
)

func Generate(err error, body any) (int, any) {
	if err == nil {
		return http.StatusOK, body
	}
	errCode, statusCode := getError(err)
	p := NewHTTPErrorPresenter()
	if statusCode >= http.StatusInternalServerError {
		fmt.Printf("[ERROR] http status: %d. %s\n", statusCode, err.Error())
	}
	return p.GenError(statusCode, errCode, err)
}

func getError(err error) (string, int) {
	if err == nil {
		return "", http.StatusOK
	}
	appErr := crooterrors.ExtractError(err)
	if appErr == nil {
		return "", http.StatusInternalServerError
	}
	errorCode := appErr.ErrorCode()
	switch errorCode {
	case crooterrors.InvalidRequestError:
		return errorCode.ToString(), http.StatusBadRequest
	case crooterrors.UnauthorizedError:
		return errorCode.ToString(), http.StatusUnauthorized
	case crooterrors.ForbiddenError:
		return errorCode.ToString(), http.StatusForbidden
	default:
		return errorCode.ToString(), http.StatusInternalServerError
	}
}
