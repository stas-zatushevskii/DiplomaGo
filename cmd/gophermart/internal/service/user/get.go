package user

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
)

func (u *ServiceUser) GetUserBalance(userID uint) (*models.UserBalance, error) {
	user, err := u.database.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &models.UserBalance{Accrual: user.CurrentBalance, WithdrawnAccrual: user.WithdrawnBalance}, nil
}
