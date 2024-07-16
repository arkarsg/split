package utils

import (
	"errors"

	decimal "github.com/greatcloak/decimal"
)

type MoneyAmount struct {
	Amount decimal.Decimal
}

func StringToMoney(money string) MoneyAmount {
	amount := Must(stringToDecimal(money))
	return MoneyAmount{
		Amount: amount,
	}
}

func stringToDecimal(money string) (decimal.Decimal, error) {
	price, err := decimal.NewFromString(money)
	return price, err
}

func (m *MoneyAmount) MoneyToString() string {
	return m.Amount.StringFixed(8)
}

func AddMoney(m1 MoneyAmount, m2 MoneyAmount) MoneyAmount {
	newAmount := m1.Amount.Add(m2.Amount)
	return MoneyAmount{
		Amount: newAmount,
	}
}

func AccumulateMonies(monies []MoneyAmount) MoneyAmount {
	total := MoneyAmount{
		Amount: decimal.Zero,
	}

	for _, money := range monies {
		total = AddMoney(money, total)
	}
	return total
}

// m1 - m2
func SubtractMoney(m1 MoneyAmount, m2 MoneyAmount) (MoneyAmount, error) {
	newAmount := m1.Amount.Sub(m2.Amount)
	if newAmount.IsNegative() {
		return MoneyAmount{}, errors.New("Subtraction will cause negative value")
	}
	return MoneyAmount{Amount: newAmount}, nil
}

func MultiplyMoney(m1 MoneyAmount, multiplier float64) MoneyAmount {
	mulFactor := decimal.NewFromFloat(multiplier)
	newAmount := m1.Amount.Mul(mulFactor)
	return MoneyAmount{
		Amount: newAmount,
	}
}
