package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"io"
	"net/http"
)

func (h *Handler) OrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "OrderCreate"
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		orderNumber := string(body)
		userID, ok := r.Context().Value(constants.UserIDKey).(uint)
		if !ok {
			http.Error(w, utils.ErrorAsJSON(customErrors.ErrUserNotFound), http.StatusUnauthorized)
		}

		err = h.service.OrderService.AddNewOrder(orderNumber, userID)
		if err != nil {
			switch {
			case errors.Is(err, customErrors.ErrOrderAlreadyExist):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusConflict)
				return
			case errors.Is(err, customErrors.ErrOrderAlreadyUsed):
				w.WriteHeader(http.StatusOK)
				return
			case errors.Is(err, customErrors.ErrOrderInvalid):
				http.Error(w, utils.ErrorAsJSON(err), http.StatusUnprocessableEntity)
				return
			default:
				h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
				http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *Handler) OrdersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "OrdersGet"
		userID, ok := r.Context().Value(constants.UserIDKey).(uint)
		if !ok {
			http.Error(w, utils.ErrorAsJSON(customErrors.ErrUserNotFound), http.StatusUnauthorized)
		}
		h.logger.Info(fmt.Sprintf(fmt.Sprintf("[%s]DOING REQUEST TO ALL ORDERS WITH --- User ID: %v", HandlerName, userID)))
		orders, err := h.service.OrderService.GetAllOrders(userID)
		h.logger.Info(fmt.Sprintf("GOT RESPONSE FROM DB -- %v", orders))
		if err != nil {
			if errors.Is(err, customErrors.ErrOrdersNotFound) {
				http.Error(w, utils.ErrorAsJSON(err), http.StatusNoContent)
			}
		}
		response, err := json.Marshal(orders)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
