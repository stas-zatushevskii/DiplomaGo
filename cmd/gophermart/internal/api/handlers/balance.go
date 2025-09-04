package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"net/http"
)

func (h *Handler) GetUserBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetUserBalance"
		userID := r.Context().Value(constants.UserIDKey).(uint)
		user, err := h.service.UserService.GetUserBalance(userID)
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
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
	Withdrawn float64 `json:"sum" validate:"required"`
}

func (h *Handler) WithdrawOrderAccrual() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "WithdrawOrderAccrual"
		var requestData withdrawBalanceData
		userID := r.Context().Value(constants.UserIDKey).(uint)
		processBodyResponse := utils.ProcessBody(&utils.ValidateData{
			R:           r,
			Logger:      h.logger,
			Validator:   h.validator,
			HandlerName: HandlerName,
			RequestData: &requestData,
		})
		if processBodyResponse.ErrCode != 0 {
			http.Error(w, processBodyResponse.ErrMsg, processBodyResponse.ErrCode)
			return
		}

		err := h.service.OrderService.AddExternalOrder(requestData.Order, userID) // adding new order in database with status Processed
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
		}
		userBalance, err := h.service.UserService.GetUserBalance(userID)
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
		}
		err = h.service.OrderService.Withdraw(
			models.ProcessOderData{UserID: userID, OrderNumber: requestData.Order},
			requestData.Withdrawn,
			userBalance.Accrual)
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetWithdrawalsHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "GetWithdrawalsHistory"
		userID := r.Context().Value(constants.UserIDKey).(uint)
		history, err := h.service.OrderService.GetWithdrawByUserID(userID)
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
		}
		response, err := json.Marshal(history)
		if err != nil {
			if err != nil {
				resp := utils.ProcessServiceError(err, h.logger, HandlerName)
				if resp.HttpStatus != http.StatusOK {
					http.Error(w, resp.ErrMsg, resp.HttpStatus)
					return
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
