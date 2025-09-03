package order

import (
	"fmt"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
)

func (o *ServiceOrder) Withdraw(orderData models.ProcessOderData, withdrawn float64, userBalance utils.Money) error {
	order, err := o.database.GetOrderByOrderNumber(orderData.OrderNumber)
	if err != nil {
		return fmt.Errorf("failed to get order by order number: %w", err)
	}
	formatedWithdrawn := utils.NewMoneyFromFloat(withdrawn)
	if userBalance <= utils.NewMoneyFromFloat(withdrawn) {
		return customErrors.ErrNotEnoughBalance
	}
	err = o.database.DecreaseUserBalance(orderData.UserID, formatedWithdrawn)
	if err != nil {
		return fmt.Errorf("error when decreasing order accrual: %w", err)
	}
	err = o.database.AddToOrderHistory(order, formatedWithdrawn)
	if err != nil {
		return fmt.Errorf("failed to add to order history: %w", err)
	}
	return nil
}

func (o *ServiceOrder) GetWithdrawByUserID(userID uint) ([]models.OrderHistory, error) {
	history, err := o.database.GetWithdrawalsHistory(userID)
	if err != nil {
		return nil, err
	}
	if len(history) == 0 {
		return nil, customErrors.ErrNoWithdrawals
	}
	return history, nil
}
