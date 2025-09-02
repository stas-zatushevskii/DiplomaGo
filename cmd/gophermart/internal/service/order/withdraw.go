package order

import (
	"fmt"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
)

func (o *ServiceOrder) Withdraw(withdrawn float64, orderNumber string) error {
	order, err := o.database.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		return fmt.Errorf("failed to get order by order number: %w", err)
	}
	formatedWithdrawn := utils.NewMoneyFromFloat(withdrawn)
	if order.Accrual <= formatedWithdrawn {
		return customErrors.ErrNotEnoughBalance
	}
	err = o.database.DecreaseOrderAccrual(orderNumber, formatedWithdrawn)
	if err != nil {
		return fmt.Errorf("error when decreasing order accrual: %w", err)
	}
	err = o.database.AddToOrderHistory(order, formatedWithdrawn)
	if err != nil {
		return fmt.Errorf("failed to add to order history: %w", err)
	}
	return nil
}

func (o *ServiceOrder) WithdrawVersion2(withdrawn float64, orderNumber string, userBalance *models.UserBalance) error {
	order, err := o.database.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		return fmt.Errorf("failed to get order by order number: %w", err)
	}
	formatedWithdrawn := utils.NewMoneyFromFloat(withdrawn)
	if userBalance.Accrual <= formatedWithdrawn {
		fmt.Println("----------------------------------------------")
		fmt.Println("Withdrawn:", formatedWithdrawn, "userBalance:", userBalance.Accrual)
		fmt.Println("----------------------------------------------")
		return customErrors.ErrNotEnoughBalance
	}
	o.logger.Warn(fmt.Sprintf("BALANCE TO WITHDRAW : %v", formatedWithdrawn))
	err = o.database.DecreaseOrderAccrualVersion2(orderNumber, formatedWithdrawn)
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
