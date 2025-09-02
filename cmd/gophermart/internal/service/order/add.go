package order

import (
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
