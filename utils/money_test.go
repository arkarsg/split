package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringToMoney(t *testing.T) {
	testString := "150.50"
	expectedMoneyAmount := MoneyAmount{
		Dollars: 150,
		Cents:   50,
	}
	moneyAmount, err := StringToMoney(testString)
	require.NoError(t, err)
	assert.Equal(t, expectedMoneyAmount, moneyAmount)
}

func TestStringToMoneyWithPaddedZeroes(t *testing.T) {
	testString := "55.00000000"
	expectedMoneyAmount := MoneyAmount{
		Dollars: 55,
		Cents:   000,
	}
	moneyAmount, err := StringToMoney(testString)
	require.NoError(t, err)
	assert.Equal(t, expectedMoneyAmount, moneyAmount)
}

func TestMoneyToStringLeadingZeroes(t *testing.T) {
	money := MoneyAmount{
		Dollars: 100,
		Cents:   0,
	}

	moneyString := money.MoneyToString()
	assert.Equal(t, "100.0", moneyString)
}

func TestMoneyToString(t *testing.T) {
	money := MoneyAmount{
		Dollars: 100,
		Cents:   5,
	}

	moneyString := money.MoneyToString()
	assert.Equal(t, "100.05", moneyString)
}

func TestAddMoney(t *testing.T) {
	money1 := MoneyAmount{
		Dollars: 0,
		Cents:   99,
	}
	money2 := MoneyAmount{
		Dollars: 0,
		Cents:   1,
	}

	expectedMoney := MoneyAmount{
		Dollars: 1,
		Cents:   0,
	}

	actualMoney := AddMoney(money1, money2)
	assert.Equal(t, expectedMoney, actualMoney)
}

func TestAccumulateMonies(t *testing.T) {
	tests := []struct {
		monies   []MoneyAmount
		expected MoneyAmount
	}{
		{[]MoneyAmount{{123, 45}, {10, 55}}, MoneyAmount{134, 0}},
		{[]MoneyAmount{{0, 99}, {0, 2}}, MoneyAmount{1, 1}},
		{[]MoneyAmount{{100, 0}, {0, 100}}, MoneyAmount{101, 0}},
	}

	for _, test := range tests {
		result := AccumulateMonies(test.monies)
		if result != test.expected {
			t.Errorf("AccumulateMonies(%v) = %v, expected %v", test.monies, result, test.expected)
		}
	}
}

func TestSubtractMoney(t *testing.T) {
	money1 := MoneyAmount{
		Dollars: 10,
		Cents:   1,
	}
	money2 := MoneyAmount{
		Dollars: 1,
		Cents:   99,
	}
	expectedMoney := MoneyAmount{
		Dollars: 8,
		Cents:   2,
	}
	actualMoney, _ := SubtractMoney(money1, money2)
	assert.Equal(t, expectedMoney, actualMoney)
}

func TestMultiplyMoney(t *testing.T) {
	money1 := MoneyAmount{
		Dollars: 10,
		Cents:   10,
	}

	expectedMoney := MoneyAmount{
		Dollars: 1,
		Cents:   1,
	}

	actualMoney := MultiplyMoney(money1, 0.1)
	assert.Equal(t, expectedMoney, actualMoney)
}

func TestMultiplyMoneyHandlesExchangeRate(t *testing.T) {
	money1 := MoneyAmount{
		Dollars: 1,
		Cents:   0,
	}

	expectedMoney := MoneyAmount{
		Dollars: 1,
		Cents:   25,
	}

	actualMoney := MultiplyMoney(money1, 1.25)
	assert.Equal(t, expectedMoney, actualMoney)
}
