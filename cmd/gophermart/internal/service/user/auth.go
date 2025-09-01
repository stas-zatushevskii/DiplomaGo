package user

import (
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"go.uber.org/zap"
)

func (u *ServiceUser) CreateNew(username, password string) (*models.User, error) {
	var user models.User
	hashedPass, err := hashPassword(password)
	if err != nil {
		u.logger.Error("Failed to hash password", zap.Error(err))
		return &user, err
	}
	user.Username = username
	user.HashedPassword = hashedPass
	err = u.database.CreateNewUser(&user)
	if err != nil {
		u.logger.Error("Failed to create new user", zap.Error(err))
		return &user, err
	}

	return &user, nil
}

type LoginResponse struct {
	UserId        uint
	Authenticated bool
}

func (u *ServiceUser) Login(username, password string) (*LoginResponse, error) {
	var response LoginResponse
	user, err := u.database.GetUserByUsername(username)
	if err != nil {
		u.logger.Info("Failed to get user", zap.Error(err))
		return &response, err
	}
	ok, err := verifyPassword(password, user.HashedPassword)

	if err != nil {
		u.logger.Error("Failed to verify password", zap.Error(err))
		return &response, err
	}
	response.UserId = user.ID
	response.Authenticated = ok
	return &response, nil
}
