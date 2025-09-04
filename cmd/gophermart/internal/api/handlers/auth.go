package handlers

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	"net/http"
	"time"
)

type UserRequest struct {
	Username string `json:"login" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum"`
}

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData UserRequest
		const HandlerName = "Register"
		processBodyResponse := utils.ProcessBody(&utils.ValidateData{
			R:           r,
			Logger:      h.logger,
			HandlerName: HandlerName,
			RequestData: requestData,
		})
		if processBodyResponse.ErrCode != 0 {
			http.Error(w, processBodyResponse.ErrMsg, processBodyResponse.ErrCode)
		}

		user, err := h.service.UserService.CreateNew(requestData.Username, requestData.Password)
		if err != nil {
			resp := utils.ProcessServiceError(err, h.logger, HandlerName)
			if resp.HttpStatus != http.StatusOK {
				http.Error(w, resp.ErrMsg, resp.HttpStatus)
				return
			}
		}
		loginResponse := h.processUserLogin(user.Username, requestData.Password, HandlerName)
		if loginResponse.HttpStatus != http.StatusOK {
			http.Error(w, loginResponse.ErrMsg, loginResponse.HttpStatus)
			return
		}

		http.SetCookie(w, loginResponse.Cookie)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData UserRequest
		const HandlerName = "Login"

		validateResponse := utils.ProcessBody(&utils.ValidateData{
			R:           r,
			Logger:      h.logger,
			HandlerName: HandlerName,
			RequestData: requestData,
		})

		if validateResponse.ErrCode != 0 {
			http.Error(w, validateResponse.ErrMsg, validateResponse.ErrCode)
		}

		loginResponse := h.processUserLogin(requestData.Username, requestData.Password, HandlerName)
		if loginResponse.HttpStatus != http.StatusOK {
			http.Error(w, loginResponse.ErrMsg, loginResponse.HttpStatus)
			return
		}

		http.SetCookie(w, loginResponse.Cookie)
		w.WriteHeader(http.StatusOK)
	}
}

type processUserLoginResponse struct {
	utils.ProcessErrorResponse
	Cookie *http.Cookie
}

// processUserLogin processing user login: Authenticate user, building jwt and creating cookie object.
// If all good -> HttpStatus = http.StatusOK, cookie = http.Cookie(...)
func (h *Handler) processUserLogin(Username, Password, HandlerName string) processUserLoginResponse {
	login, err := h.service.UserService.Login(Username, Password)
	if err != nil {
		return processUserLoginResponse{
			ProcessErrorResponse: utils.ProcessServiceError(err, h.logger, HandlerName),
			Cookie:               nil,
		}
	}
	if !login.Authenticated {
		return processUserLoginResponse{
			ProcessErrorResponse: utils.ProcessErrorResponse{HttpStatus: http.StatusUnauthorized, ErrMsg: "unauthorized"},
			Cookie:               nil,
		}
	}
	jwt, err := h.service.UserService.BuildJwt(login.UserID)
	if err != nil {
		return processUserLoginResponse{
			ProcessErrorResponse: utils.ProcessServiceError(err, h.logger, HandlerName),
			Cookie:               nil,
		}
	}
	cookie := &http.Cookie{
		Name:    "JWT",
		Expires: time.Now().Add(24 * time.Hour),
		Value:   jwt,
	}
	return processUserLoginResponse{
		ProcessErrorResponse: utils.ProcessErrorResponse{HttpStatus: http.StatusOK, ErrMsg: ""},
		Cookie:               cookie,
	}
}
