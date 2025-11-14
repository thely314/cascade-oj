package util

import (
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	stored, err := GenerateHashPassword("password")
	if err != nil {
		t.Fatalf("GenerateHashPassword failed: %v", err)
	}
	tests := []struct {
		password string
		stored   string
		expected bool
	}{
		{"password", stored,
			true},
		{"wrongpassword", stored,
			false},
	}
	for _, tt := range tests {
		result := VerifyPassword(tt.password, tt.stored)
		if result != tt.expected {
			t.Errorf("VerifyPassword(%q, %q) = %v; want %v",
				tt.password, tt.stored, result, tt.expected)
		}
	}
}
