package order

import (
	"fmt"
	customErr "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
)

func (o *ServiceOrder) AddNewOrder(orderNumber string, userID uint, orderChan chan<- models.ProcessOderData) error {
	ok := CheckLuhna(orderNumber)
	if !ok {
		return customErr.ErrOrderInvalid
	}
	_, err := o.database.CreateNewOrder(orderNumber, userID)

	if err == nil {
		go func() {
			orderChan <- models.ProcessOderData{UserID: userID, OrderNumber: orderNumber}
		}()
	}
	return err
}

func (o *ServiceOrder) AddExternalOrder(orderNumber string, userID uint) error {
	ok := CheckLuhna(orderNumber)
	if !ok {
		return customErr.ErrOrderInvalid
	}
	_, err := o.database.CreateProcessedOrder(orderNumber, userID)
	if err != nil {
		return fmt.Errorf("error adding processed order: %w", err)
	}
	return nil
}
