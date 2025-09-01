package order

import (
	customErr "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
)

func (o *ServiceOrder) AddNewOrder(orderNumber string, userID uint) error {
	ok := CheckLuhna(orderNumber)
	if !ok {
		return customErr.ErrOrderInvalid
	}
	_, err := o.database.CreateNewOrder(orderNumber, userID)
	return err
}
