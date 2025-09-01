package order

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
)

func (o *ServiceOrder) GetAllOrders(userID uint) ([]models.Order, error) {
	var response []models.Order
	orders, err := o.database.GetOrdersByUserId(userID)
	if err != nil {
		return response, err
	}
	return orders, nil
}
