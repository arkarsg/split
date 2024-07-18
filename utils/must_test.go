package utils

import (
	"errors"
	"testing"
)

func TestMust_NoError(t *testing.T) {
	expected := 42
	result := Must(expected, nil)
	if result != expected {
		t.Errorf("Must() = %v, want %v", result, expected)
	}
}

func TestMust_WithError(t *testing.T) {
	expectedErr := errors.New("some error")
	defer func() {
		if r := recover(); r != nil {
			if r != expectedErr {
				t.Errorf("Must() panic = %v, want %v", r, expectedErr)
			}
		} else {
			t.Errorf("Must() did not panic")
		}
	}()
	Must(0, expectedErr)
}
