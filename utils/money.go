package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type MoneyAmount struct {
	Dollars uint64
	Cents   uint64
}

func StringToMoney(money string) MoneyAmount {
	return Must(stringToMoney(money))
}

func stringToMoney(money string) (MoneyAmount, error) {
	dollarsAndCents := strings.Split(money, ".")
	if len(dollarsAndCents) > 2 {
		return MoneyAmount{}, errors.New("There are more than 2 elements to dollars and cents")
	}

	dollars, err := strconv.Atoi(dollarsAndCents[0])
	if err != nil {
		return MoneyAmount{}, errors.New("Error converting dollar to int")
	}
	cents, err := strconv.Atoi(dollarsAndCents[1])
	if err != nil || cents > 99 {
		return MoneyAmount{}, errors.New("Error converting cents to int")
	}
	return MoneyAmount{
		Dollars: uint64(dollars),
		Cents:   uint64(cents),
	}, nil
}

func (m *MoneyAmount) MoneyToString() string {
	dollars := strconv.Itoa(int(m.Dollars))
	cents := strconv.Itoa(int(m.Cents))
	if m.Cents < 10 {
		cents = "0" + cents
	}
	return dollars + "." + cents
}

func AddMoney(m1 MoneyAmount, m2 MoneyAmount) MoneyAmount {
	totalCents := m1.Cents + m2.Cents
	extraDollars := totalCents / 100
	totalDollars := m1.Dollars + m2.Dollars + uint64(extraDollars)
	remainingCents := totalCents % 100

	return MoneyAmount{
		Dollars: totalDollars,
		Cents:   remainingCents,
	}
}

func AccumulateMonies(monies []MoneyAmount) MoneyAmount {
	total := MoneyAmount{
		Dollars: 0,
		Cents:   0,
	}

	for _, money := range monies {
		total = AddMoney(money, total)
	}
	return total
}

func SubtractMoney(m1 MoneyAmount, m2 MoneyAmount) (MoneyAmount, error) {
	if m1.Dollars < m2.Dollars || (m1.Dollars == m2.Dollars && m1.Cents < m2.Cents) {
		return MoneyAmount{}, errors.New("Subtraction result would be negative")
	}

	totalCents1 := m1.Dollars*100 + uint64(m1.Cents)
	totalCents2 := m2.Dollars*100 + uint64(m2.Cents)

	if totalCents1 < totalCents2 {
		return MoneyAmount{}, errors.New("Subtraction result would be negative")
	}

	remainingCents := totalCents1 - totalCents2

	return MoneyAmount{
		Dollars: remainingCents / 100,
		Cents:   uint64(remainingCents % 100),
	}, nil
}

func MultiplyMoney(m1 MoneyAmount, multiplier float64) MoneyAmount {
	totalCents := float64(m1.Dollars*100+uint64(m1.Cents)) * multiplier
	dollars := uint64(totalCents) / 100
	cents := uint64(math.Round(totalCents)) % 100

	return MoneyAmount{
		Dollars: dollars,
		Cents:   cents,
	}
}
