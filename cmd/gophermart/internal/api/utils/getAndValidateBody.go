package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type ValidateData struct {
	R           *http.Request
	Logger      *zap.Logger
	Validator   *validator.Validate
	HandlerName string
	RequestData interface{}
}

type ValidateResponse struct {
	ErrMsg  string
	ErrCode int
}

// ProcessBody getting data from request body, validating json structure, unmarshalling json, validating json data.
// If OK -> errStr = "", errCode = 0
func ProcessBody(data *ValidateData) (response ValidateResponse) {
	body, err := io.ReadAll(data.R.Body)
	if err != nil {
		return ValidateResponse{
			ErrMsg:  ErrorAsJSON(err),
			ErrCode: http.StatusInternalServerError,
		}
	}
	// safe Body.Close()
	defer func() {
		if err := data.R.Body.Close(); err != nil {
			data.Logger.Error(fmt.Sprintf("%s: error closing body", data.HandlerName), zap.Error(err))
		}
	}()

	if !json.Valid(body) {
		return ValidateResponse{
			ErrMsg:  "Invalid json",
			ErrCode: http.StatusBadRequest,
		}
	}

	err = json.Unmarshal(body, data.RequestData)
	if err != nil {
		return ValidateResponse{
			ErrMsg:  ErrorAsJSON(err),
			ErrCode: http.StatusInternalServerError,
		}
	}

	// Validate by tags
	err = data.Validator.Struct(data.RequestData)
	if err != nil {
		return ValidateResponse{
			ErrMsg:  ErrorAsJSON(err),
			ErrCode: http.StatusBadRequest,
		}
	}
	return
}
