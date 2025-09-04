package errors

import "errors"

// Сигнальные ошибки см. 100 Ошибок Go стр. 220
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrOrdersNotFound         = errors.New("orders not found")
	ErrOrderInvalid           = errors.New("order invalid")
	ErrOrderAlreadyExist      = errors.New("order already exist")
	ErrOrderAlreadyUsed       = errors.New("order already used")
	ErrNotEnoughBalance       = errors.New("not enough balance")
	ErrNoWithdrawals          = errors.New("no withdrawals")
	ErrConfigNotFound         = errors.New("config not found")
	ErrUserNotFoundByUsername = errors.New("user not found by username")
	ErrJWTBuild               = errors.New("failed to build JWT token")
)
