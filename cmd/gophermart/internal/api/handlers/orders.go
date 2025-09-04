package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/api/utils"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"net/http"
)

func (h *Handler) OrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "OrderCreate"
		userID := r.Context().Value(constants.UserIDKey).(uint)

		orderNumber, err := utils.GetTextPlain(r, h.logger, HandlerName)
		if err != nil {
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}

		err = h.service.OrderService.AddNewOrder(orderNumber, userID, h.orderChan)
		if err != nil {
			resp := utils.ProcessServiceError(err, h.logger, HandlerName)
			if resp.HttpStatus != http.StatusOK {
				http.Error(w, resp.ErrMsg, resp.HttpStatus)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) OrdersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const HandlerName = "OrdersGet"
		userID := r.Context().Value(constants.UserIDKey).(uint)

		orders, err := h.service.OrderService.GetAllOrders(userID)
		if err != nil {
			resp := utils.ProcessServiceError(err, h.logger, HandlerName)
			if resp.HttpStatus != http.StatusOK {
				http.Error(w, resp.ErrMsg, resp.HttpStatus)
				return
			}
		}
		response, err := json.Marshal(orders)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%s: %s", HandlerName, err.Error()))
			http.Error(w, utils.ErrorAsJSON(err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
