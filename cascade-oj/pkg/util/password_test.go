package util

import (
	"crypto/md5"
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	md5Password := md5.Sum([]byte("password"))
	md5PasswordStr := string(md5Password[:])
	stored, err := GenerateHashPassword(md5PasswordStr)
	if err != nil {
		t.Fatalf("GenerateHashPassword failed: %v", err)
	}
	tests := []struct {
		password string
		stored   string
		expected bool
	}{
		{md5PasswordStr, stored,
			true},
		{"password", stored,
			false},
		{"wrongpassword", stored,
			false},
		{"", stored,
			false},
		{md5PasswordStr, "",
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
