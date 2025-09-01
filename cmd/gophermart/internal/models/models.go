package models

import (
	"encoding/json"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model     `json:"-"`
	Username       string `gorm:"unique" json:"-"`
	HashedPassword string `json:"-"`
	Orders         []Order
}

type UserBalance struct {
	Accrual          utils.Money `json:"accrual"`
	WithdrawnAccrual utils.Money `json:"withdrawn"`
}

func (u UserBalance) MarshalJSON() (data []byte, err error) {
	type aliasData UserBalance
	aliasValue := struct {
		aliasData
		Accrual          float64 `json:"accrual"`
		WithdrawnAccrual float64 `json:"withdrawn"`
	}{
		aliasData:        aliasData(u),
		Accrual:          utils.Money.ToFloat(u.Accrual),
		WithdrawnAccrual: utils.Money.ToFloat(u.WithdrawnAccrual),
	}
	return json.Marshal(aliasValue)
}

type Order struct {
	ID               uint                  `gorm:"primaryKey" json:"-"`
	CreatedAt        string                `json:"uploaded_at"`
	OrderNumber      string                `gorm:"uniqueIndex;not null" json:"number"`
	UserID           *uint                 `json:"-"`
	User             User                  `json:"-"`
	Status           constants.OrderStatus `json:"status"`
	Accrual          utils.Money           `json:"accrual,omitempty"`
	WithdrawnAccrual utils.Money           `json:"withdrawn,omitempty"`
	History          []OrderHistory
}

func (u Order) MarshalJSON() (data []byte, err error) {
	type aliasData Order
	aliasValue := struct {
		aliasData
		Accrual          float64 `json:"accrual,omitempty"`
		WithdrawnAccrual float64 `json:"withdrawn,omitempty"`
	}{
		aliasData:        aliasData(u),
		Accrual:          utils.Money.ToFloat(u.Accrual),
		WithdrawnAccrual: utils.Money.ToFloat(u.WithdrawnAccrual),
	}
	return json.Marshal(aliasValue)
}

type OrderHistory struct {
	gorm.Model  `json:"-"`
	ID          uint        `gorm:"primaryKey" json:"-"`
	OrderID     *uint       `gorm:"uniqueIndex" json:"-"`
	OrderNumber string      `json:"order"`
	Sum         utils.Money `json:"sum"`
	ProcessedAt string      `json:"processed_at"`
	Order       Order       `json:"-"`
}

func (u OrderHistory) MarshalJSON() (data []byte, err error) {
	type aliasData OrderHistory
	aliasValue := struct {
		aliasData
		Sum float64 `json:"accrual"`
	}{
		aliasData: aliasData(u),
		Sum:       utils.Money.ToFloat(u.Sum),
	}
	return json.Marshal(aliasValue)
}

type AccrualResponse struct {
	Status  constants.OrderStatus `json:"status"`
	Accrual utils.Money           `json:"accrual"`
	Order   string                `json:"order"`
}

func (u AccrualResponse) MarshalJSON() (data []byte, err error) {
	type aliasData AccrualResponse
	aliasValue := struct {
		aliasData
		Accrual float64 `json:"accrual"`
	}{
		aliasData: aliasData(u),
		Accrual:   utils.Money.ToFloat(u.Accrual),
	}
	return json.Marshal(aliasValue)
}

type ProcessOrderData struct {
	Accrual     utils.Money `json:"accrual"`
	OrderNumber string      `json:"order_number"`
}

func (u ProcessOrderData) MarshalJSON() (data []byte, err error) {
	type aliasData ProcessOrderData
	aliasValue := struct {
		aliasData
		Accrual float64 `json:"accrual"`
	}{
		aliasData: aliasData(u),
		Accrual:   utils.Money.ToFloat(u.Accrual),
	}
	return json.Marshal(aliasValue)
}
