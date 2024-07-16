package utils

import (
	"testing"

	"github.com/greatcloak/decimal"
	"github.com/stretchr/testify/assert"
)

func TestStringToMoney(t *testing.T) {
	testString := "150.50000000"
	moneyAmount := StringToMoney(testString)
	assert.Equal(t, testString, moneyAmount.MoneyToString())
}

func TestStringToMoneyWithPaddedZeroes(t *testing.T) {
	testString := "55.00000001"
	moneyAmount := StringToMoney(testString)
	assert.Equal(t, testString, moneyAmount.MoneyToString())
}

func TestMoneyToString(t *testing.T) {
	money := MoneyAmount{
		Amount: decimal.NewFromFloat(10.05),
	}

	moneyString := money.MoneyToString()
	assert.Equal(t, "10.05000000", moneyString)
}

func TestAddMoney(t *testing.T) {
	money1 := MoneyAmount{
		Amount: decimal.NewFromFloat(0.99),
	}

	money2 := MoneyAmount{
		Amount: decimal.NewFromFloat(0.01),
	}

	expectedMoneyString := "1.00000000"

	actualMoney := AddMoney(money1, money2)
	assert.Equal(t, expectedMoneyString, actualMoney.MoneyToString())
}

func TestAccumulateMonies(t *testing.T) {
	tests := []struct {
		name   string
		monies []MoneyAmount
		expect MoneyAmount
	}{
		{
			name:   "Empty list",
			monies: []MoneyAmount{},
			expect: MoneyAmount{Amount: decimal.Zero},
		},
		{
			name: "Single money",
			monies: []MoneyAmount{
				{Amount: decimal.NewFromFloat(10.50)},
			},
			expect: MoneyAmount{Amount: decimal.NewFromFloat(10.50)},
		},
		{
			name: "Multiple monies",
			monies: []MoneyAmount{
				{Amount: decimal.NewFromFloat(5.25)},
				{Amount: decimal.NewFromFloat(3.75)},
				{Amount: decimal.NewFromFloat(1.50)},
			},
			expect: MoneyAmount{Amount: decimal.NewFromFloat(10.50)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AccumulateMonies(tt.monies)
			assert.True(t, tt.expect.Amount.Equal(result.Amount))
		})
	}
}

func TestSubtractMoney(t *testing.T) {
	tests := []struct {
		name     string
		m1       MoneyAmount
		m2       MoneyAmount
		expected MoneyAmount
		errMsg   string
	}{
		{
			name:     "Valid subtraction",
			m1:       MoneyAmount{Amount: decimal.NewFromFloat(10.00)},
			m2:       MoneyAmount{Amount: decimal.NewFromFloat(3.50)},
			expected: MoneyAmount{Amount: decimal.NewFromFloat(6.50)},
			errMsg:   "",
		},
		{
			name:     "Subtraction resulting in negative",
			m1:       MoneyAmount{Amount: decimal.NewFromFloat(3.00)},
			m2:       MoneyAmount{Amount: decimal.NewFromFloat(5.00)},
			expected: MoneyAmount{},
			errMsg:   "Subtraction will cause negative value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SubtractMoney(tt.m1, tt.m2)
			if tt.errMsg != "" {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.expected.Amount.Equal(result.Amount))
			}
		})
	}
}

func TestMultiplyMoney(t *testing.T) {
	tests := []struct {
		name       string
		m1         MoneyAmount
		multiplier float64
		expected   MoneyAmount
	}{
		{
			name:       "Multiply by 2",
			m1:         MoneyAmount{Amount: decimal.NewFromFloat(5.00)},
			multiplier: 2,
			expected:   MoneyAmount{Amount: decimal.NewFromFloat(10.00)},
		},
		{
			name:       "Multiply by 0.5",
			m1:         MoneyAmount{Amount: decimal.NewFromFloat(10.00)},
			multiplier: 0.5,
			expected:   MoneyAmount{Amount: decimal.NewFromFloat(5.00)},
		},
		{
			name:       "Multiply by 1.5",
			m1:         MoneyAmount{Amount: decimal.NewFromFloat(7.50)},
			multiplier: 1.5,
			expected:   MoneyAmount{Amount: decimal.NewFromFloat(11.25)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MultiplyMoney(tt.m1, tt.multiplier)
			assert.True(t, tt.expected.Amount.Equal(result.Amount))
		})
	}
}
