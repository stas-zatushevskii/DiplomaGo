package utils

import (
	"errors"
	"fmt"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"go.uber.org/zap"
	"net/http"
)

type ProcessErrorResponse struct {
	ErrMsg     string
	HttpStatus int
}

// ProcessServiceError return httpStatus (200, 400 ...) and errMsg according to custom error type.
func ProcessServiceError(err error, Logger *zap.Logger, HandlerName string) ProcessErrorResponse {
	switch {
	case errors.Is(err, customErrors.ErrUserNotFoundByUsername):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusNotFound,
		}
	case errors.Is(err, customErrors.ErrUserNotFound):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusUnauthorized,
		}
	case errors.Is(err, customErrors.ErrOrderAlreadyExist):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusConflict,
		}
	case errors.Is(err, customErrors.ErrOrderAlreadyUsed):
		return ProcessErrorResponse{
			ErrMsg:     "",
			HttpStatus: http.StatusOK,
		}
	case errors.Is(err, customErrors.ErrOrderInvalid):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusUnprocessableEntity,
		}
	case errors.Is(err, customErrors.ErrOrdersNotFound):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusUnprocessableEntity,
		}
	case errors.Is(err, customErrors.ErrNotEnoughBalance):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusPaymentRequired,
		}
	case errors.Is(err, customErrors.ErrNoWithdrawals):
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusNoContent,
		}
	default:
		Logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
		return ProcessErrorResponse{
			ErrMsg:     ErrorAsJSON(err),
			HttpStatus: http.StatusInternalServerError,
		}
	}
}
