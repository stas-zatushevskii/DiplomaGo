package utils

import (
	"math"
)

type Money int64

func NewMoneyFromFloat(v float64) Money {
	return Money(math.Round(v * 100))
}

func (m Money) ToFloat() float64 {
	return float64(m) / 100
}
