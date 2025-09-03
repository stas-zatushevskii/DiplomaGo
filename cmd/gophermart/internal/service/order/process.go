package order

import (
	"context"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
	"go.uber.org/zap"
	"time"
)

var ProcessingOrdersCache = make(map[string]struct{})

func (o *ServiceOrder) ProcessOrder(data models.ProcessOderData) error { // user models.User
	ProcessingOrdersCache[data.OrderNumber] = struct{}{}
	defer delete(ProcessingOrdersCache, data.OrderNumber)

	// send request with retry
	accrualResponse, err := o.RequestWithRetry(o, data.OrderNumber)
	if err != nil {
		return err
	}
	switch accrualResponse.Status {
	case constants.OrderStatusProcessing:
		err = o.database.ChangeOrderStatus(data.OrderNumber, constants.OrderStatusProcessing)
		if err != nil {
			return err
		}
	case constants.OrderStatusInvalid:
		err = o.database.ChangeOrderStatus(data.OrderNumber, constants.OrderStatusInvalid)
		if err != nil {
			return err
		}
	case constants.OrderStatusProcessed:
		err = o.database.IncreaseOrderAccrual(data.OrderNumber, accrualResponse.Accrual)
		if err != nil {
			return err
		}
		err = o.database.IncreaseUserBalance(data.UserID, accrualResponse.Accrual)
		if err != nil {
			return err
		}
		err = o.database.ChangeOrderStatus(data.OrderNumber, constants.OrderStatusProcessed)
		if err != nil {
			return err
		}
	}
	return nil
}

// OrderListener load orders from orderChan and process them
func (o *ServiceOrder) OrderListener(ctx context.Context, ch <-chan models.ProcessOderData, semaphore *utils.Semaphore) {
	for {
		select {
		case <-ctx.Done():
			o.logger.Info("order listener stopped", zap.Error(ctx.Err()))
			return
		case orderData := <-ch:
			semaphore.Acquire()
			go func(orderNumber models.ProcessOderData) {
				defer semaphore.Release()
				err := o.ProcessOrder(orderData)
				if err != nil {
					o.logger.Error("process order failed", zap.Error(err))
				}
			}(orderData)
		}
	}
}

func (o *ServiceOrder) loadOrders(ch chan<- models.ProcessOderData) {
	data, err := o.database.GetNewOrProcessingOrders()
	if err != nil {
		o.logger.Error("get new or processing orders failed", zap.Error(err))
		return
	}
	for _, order := range data {
		_, ok := ProcessingOrdersCache[order.OrderNumber]
		if !ok {
			o.logger.Info("processing order from database " + order.OrderNumber)
			ch <- models.ProcessOderData{UserID: *order.UserID, OrderNumber: order.OrderNumber}
		}
	}
}

// OrderLoader loads orders with status New or Processing from database and sends to orderChan
func (o *ServiceOrder) OrderLoader(ctx context.Context, ch chan<- models.ProcessOderData) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// init load
	o.loadOrders(ch)

	for {
		select {
		case <-ctx.Done():
			o.logger.Info("order loader stopped", zap.Error(ctx.Err()))
			return
		case <-ticker.C:
			o.loadOrders(ch)
		}
	}
}
