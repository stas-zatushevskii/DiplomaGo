package utils

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// GetTextPlain getting string from request where content type = text/plain
func GetTextPlain(r *http.Request, Logger *zap.Logger, HandlerName string) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
		return "", err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			Logger.Error(fmt.Sprintf("%s: error closing body", HandlerName), zap.Error(err))
		}
	}()

	orderNumber := string(body)
	return orderNumber, nil
}
