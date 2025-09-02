package models

import (
	"encoding/json"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model       `json:"-"`
	Username         string      `gorm:"unique" json:"-"`
	HashedPassword   string      `json:"-"`
	CurrentBalance   utils.Money `json:"Current"`
	WithdrawnBalance utils.Money `json:"Withdrawn"`
	Orders           []Order
}

// FIXME UserBalance - old version

type UserBalance struct {
	Accrual          utils.Money `json:"Current"`
	WithdrawnAccrual utils.Money `json:"Withdrawn"`
}

type ProcessOderData struct {
	UserID      uint   `json:"user_id"`
	OrderNumber string `json:"order_number"`
}

func (u UserBalance) MarshalJSON() (data []byte, err error) {
	type aliasData UserBalance
	aliasValue := struct {
		aliasData
		Accrual          float64 `json:"Current"`
		WithdrawnAccrual float64 `json:"Withdrawn"`
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

func (u *AccrualResponse) UnmarshalJSON(bytes []byte) error {
	type tmp struct {
		Status  constants.OrderStatus `json:"status"`
		Accrual float64               `json:"accrual"`
		Order   string                `json:"order"`
	}
	var t tmp
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return err
	}
	u.Status = t.Status
	u.Order = t.Order
	u.Accrual = utils.NewMoneyFromFloat(t.Accrual)
	return nil
}
