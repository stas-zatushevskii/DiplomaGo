package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	CustomErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"io"
	"net/http"
)

func (h *Handler) GetUserBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetUserBalance"
		userID, ok := r.Context().Value(constants.UserIDKey).(uint)
		if !ok {
			http.Error(w, utils.ErrorAsJSON(CustomErrors.ErrUserNotFound), http.StatusUnauthorized)
		}
		user, err := h.service.UserService.GetUserBalance(userID)
		h.logger.Info(fmt.Sprintf("USER: %v", user))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(user)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

type withdrawBalanceData struct {
	Order     string  `json:"order" validate:"required"`
	Withdrawn float64 `json:"withdrawn" validate:"required"`
}

func (h *Handler) WithdrawOrderAccrual() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "WithdrawOrderAccrual"
		var requestData withdrawBalanceData
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
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusBadRequest)
			return
		}
		err = h.service.OrderService.Withdraw(requestData.Withdrawn, requestData.Order)
		if err != nil {
			switch {
			case errors.Is(err, CustomErrors.ErrOrdersNotFound):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusUnprocessableEntity)
				return
			case errors.Is(err, CustomErrors.ErrNotEnoughBalance):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusPaymentRequired)
				return
			default:
				h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
				http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetWithdrawalsHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetWithdrawalsHistory"
		userID, ok := r.Context().Value(constants.UserIDKey).(uint)
		if !ok {
			http.Error(w, utils.ErrorAsJSON(CustomErrors.ErrUserNotFound), http.StatusUnauthorized)
		}
		history, err := h.service.OrderService.GetWithdrawByUserID(userID)
		if err != nil {
			switch {
			case errors.Is(err, CustomErrors.ErrNoWithdrawals):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusNoContent)
				return
			default:
				http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
				return
			}
		}
		response, err := json.Marshal(history)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
