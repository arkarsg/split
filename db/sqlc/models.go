// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencySGD Currency = "SGD"
)

func (e *Currency) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Currency(s)
	case string:
		*e = Currency(s)
	default:
		return fmt.Errorf("unsupported scan type for Currency: %T", src)
	}
	return nil
}

type NullCurrency struct {
	Currency Currency `json:"currency"`
	Valid    bool     `json:"valid"` // Valid is true if Currency is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCurrency) Scan(value interface{}) error {
	if value == nil {
		ns.Currency, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Currency.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCurrency) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Currency), nil
}

type Debt struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	SettledAmount string `json:"settled_amount"`
}

type DebtDebtor struct {
	DebtID   int64    `json:"debt_id"`
	DebtorID int64    `json:"debtor_id"`
	Amount   string   `json:"amount"`
	Currency Currency `json:"currency"`
}

type Payment struct {
	ID        int64     `json:"id"`
	DebtID    int64     `json:"debt_id"`
	DebtorID  int64     `json:"debtor_id"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	Currency  Currency  `json:"currency"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	Amount    string    `json:"amount"`
	Currency  Currency  `json:"currency"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	PayerID   int64     `json:"payer_id"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
