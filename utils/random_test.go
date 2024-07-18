package utils

import (
	"regexp"
	"strconv"
	"testing"
)

// TestRandomInt checks if the generated random integer is within the specified range
func TestRandomInt(t *testing.T) {
	min, max := int64(10), int64(20)
	for i := 0; i < 100; i++ {
		r := RandomInt(min, max)
		if r < min || r > max {
			t.Errorf("RandomInt() = %v, want between %v and %v", r, min, max)
		}
	}
}

// TestRandomString checks if the generated string has the correct length and only contains alphabets
func TestRandomString(t *testing.T) {
	n := 10
	for i := 0; i < 100; i++ {
		r := RandomString(n)
		if len(r) != n {
			t.Errorf("RandomString() = %v, want length %v", r, n)
		}
		matched, _ := regexp.MatchString("^[a-z]+$", r)
		if !matched {
			t.Errorf("RandomString() = %v, want only alphabets", r)
		}
	}
}

// TestRandomUser checks if the generated username has the correct length
func TestRandomUser(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomUser()
		if len(r) != 6 {
			t.Errorf("RandomUser() = %v, want length %v", r, 6)
		}
	}
}

// TestRandomEmail checks if the generated email has the correct format
func TestRandomEmail(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomEmail()
		if len(r) < 7 || r[len(r)-9:] != "@test.com" {
			t.Errorf("RandomEmail() = %v, want format %v", r, "prefix@test.com")
		}
	}
}

// TestRandomAmount checks if the generated amount is within the specified range and has the correct format
func TestRandomAmount(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomAmount()
		if len(r) < 5 || r[len(r)-3:] != ".00" {
			t.Errorf("RandomAmount() = %v, want format %v", r, "amount.00")
		}
		amt, err := strconv.ParseInt(r[:len(r)-3], 10, 64)
		if err != nil {
			t.Errorf("RandomAmount() parsing error: %v", err)
		}
		if amt < 100 || amt > 1000 {
			t.Errorf("RandomAmount() = %v, want between %v and %v", amt, 100, 1000)
		}
	}
}

// TestRandomSmallAmount checks if the generated small amount is within the specified range and has the correct format
func TestRandomSmallAmount(t *testing.T) {
	for i := 0; i < 100; i++ {
		r := RandomSmallAmount()
		if len(r) < 4 || r[len(r)-3:] != ".00" {
			t.Errorf("RandomSmallAmount() = %v, want format %v", r, "amount.00")
		}
		amt, err := strconv.ParseInt(r[:len(r)-3], 10, 64)
		if err != nil {
			t.Errorf("RandomSmallAmount() parsing error: %v", err)
		}
		if amt < 1 || amt > 10 {
			t.Errorf("RandomSmallAmount() = %v, want between %v and %v", amt, 1, 10)
		}
	}
}
