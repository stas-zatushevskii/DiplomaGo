package order

import (
	"context"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
	"go.uber.org/zap"
	"time"
)

var ProcessingOrdersCache = make(map[string]struct{})

func (o *ServiceOrder) ProcessOrder(orderNumber string) error {
	ProcessingOrdersCache[orderNumber] = struct{}{}
	defer delete(ProcessingOrdersCache, orderNumber)

	// send request with retry
	accrualResponse, err := o.RequestWithRetry(o, orderNumber)
	if err != nil {
		return err
	}
	switch accrualResponse.Status {
	case constants.OrderStatusProcessing:
		err := o.database.ChangeOrderStatus(orderNumber, constants.OrderStatusProcessing)
		if err != nil {
			return err
		}
	case constants.OrderStatusInvalid:
		err := o.database.ChangeOrderStatus(orderNumber, constants.OrderStatusInvalid)
		if err != nil {
			return err
		}
	case constants.OrderStatusProcessed:
		err := o.database.IncreaseOrderAccrual(orderNumber, accrualResponse.Accrual)
		if err != nil {
			return err
		}
		err = o.database.ChangeOrderStatus(orderNumber, constants.OrderStatusProcessed)
		if err != nil {
			return err
		}
	}
	return nil
}

// OrderListener load orders from orderChan and process them
func (o *ServiceOrder) OrderListener(ctx context.Context, ch <-chan string, semaphore *utils.Semaphore) {
	for {
		select {
		case <-ctx.Done():
			o.logger.Info("order listener stopped", zap.Error(ctx.Err()))
			return
		case orderNumber := <-ch:
			semaphore.Acquire()
			go func(orderNumber string) {
				defer semaphore.Release()
				err := o.ProcessOrder(orderNumber)
				if err != nil {
					o.logger.Error("process order failed", zap.Error(err))
				}
			}(orderNumber)
		}
	}
}

func (o *ServiceOrder) loadOrders(ch chan<- string) {
	data, err := o.database.GetNewOrProcessingOrders()
	if err != nil {
		o.logger.Error("get new or processing orders failed", zap.Error(err))
		return
	}
	for _, order := range data {
		_, ok := ProcessingOrdersCache[order.OrderNumber]
		if !ok {
			o.logger.Info("processing order from database " + order.OrderNumber)
			ch <- order.OrderNumber
		}
	}
}

// OrderLoader loads orders with status New or Processing from database and sends to orderChan
func (o *ServiceOrder) OrderLoader(ctx context.Context, ch chan<- string) {
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
