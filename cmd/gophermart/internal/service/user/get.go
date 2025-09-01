package user

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
)

func (u *ServiceUser) GetUser(userID uint) (*models.User, error) {
	var response *models.User
	user, err := u.database.GetUserById(userID)
	if err != nil {
		return response, err
	}
	return user, nil
}

func (u *ServiceUser) GetUserBalance(userID uint) (*models.UserBalance, error) {

	accrual, err := u.database.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}
	withdrawnAccrual, err := u.database.GetUserWithDrawnBalance(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserBalance{Accrual: accrual, WithdrawnAccrual: withdrawnAccrual}, nil
}
