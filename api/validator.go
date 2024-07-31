package api

import (
	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/arkarsg/splitapp/utils"
	"github.com/go-playground/validator/v10"
)

var moneyAmount validator.Func = func(fl validator.FieldLevel) bool {
	amount, ok := fl.Field().Interface().(string)
	if ok {
		return utils.IsValidAmount(amount)
	}
	return false
}

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	curr, ok := fl.Field().Interface().(db.Currency)
	if ok {
		return curr == db.CurrencySGD || curr == db.CurrencyUSD
	}
	return false
}
