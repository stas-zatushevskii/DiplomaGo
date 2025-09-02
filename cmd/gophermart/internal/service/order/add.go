package order

import (
	"fmt"
	customErr "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
)

func (o *ServiceOrder) AddNewOrder(orderNumber string, userID uint, orderChan chan<- string) error {
	ok := CheckLuhna(orderNumber)
	if !ok {
		return customErr.ErrOrderInvalid
	}
	_, err := o.database.CreateNewOrder(orderNumber, userID)

	if err == nil {
		go func() {
			orderChan <- orderNumber
		}()
	}
	return err
}

func (o *ServiceOrder) AddNewSingleOrder(orderNumber string, userID uint) error {
	ok := CheckLuhna(orderNumber)
	if !ok {
		return customErr.ErrOrderInvalid
	}
	newOrder, err := o.database.CreateNewOrder(orderNumber, userID)
	if err != nil {
		return fmt.Errorf("error adding single order: %w", err)
	}
	o.logger.Info(fmt.Sprintf("Adding single order to the order list: %s", newOrder.OrderNumber))
	err = o.ProcessOrder(newOrder.OrderNumber)
	if err != nil {
		return fmt.Errorf("error processing single order: %w", err)
	}

	return nil
}
