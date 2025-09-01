package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"io"
	"net/http"
	"time"
)

type UserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum"`
}

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData UserRequest
		const HandlerName = "Login"
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		if !json.Valid(body) {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}

		// Validate by tags
		err = h.validator.Struct(requestData)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusBadRequest)
			return
		}

		user, err := h.service.UserService.CreateNew(requestData.Username, requestData.Password)
		if err != nil {
			switch {
			case errors.Is(err, customErrors.ErrUserAlreadyExists):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusConflict)
			default:
				http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			}
			return
		}
		login, err := h.service.UserService.Login(user.Username, requestData.Password)

		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		if !login.Authenticated {
			http.Error(w, utils.ErrorAsJSON(fmt.Errorf("unauthorized")), http.StatusUnauthorized)
			return
		}
		jwt, err := h.service.UserService.BuildJwt(login.UserID)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "JWT",
			Expires: time.Now().Add(24 * time.Hour),
			Value:   jwt,
		})
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData UserRequest
		const HandlerName = "Login"
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		if !json.Valid(body) {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}

		// Validate by tags
		err = h.validator.Struct(requestData)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusBadRequest)
			return
		}
		login, err := h.service.UserService.Login(requestData.Username, requestData.Password)
		if err != nil {
			if errors.Is(err, customErrors.ErrUserNotFound) {
				http.Error(w, utils.ErrorAsJSON(err), http.StatusNotFound)
				return
			}
		}
		if !login.Authenticated {
			http.Error(w, utils.ErrorAsJSON(fmt.Errorf("unauthorized")), http.StatusUnauthorized)
			return
		}
		jwt, err := h.service.UserService.BuildJwt(login.UserID)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "JWT",
			Expires: time.Now().Add(24 * time.Hour),
			Value:   jwt,
		})
		w.WriteHeader(http.StatusOK)
	}
}
