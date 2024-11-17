package jwt

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(JwtPayload{
		UserID:   1,
		Role:     1,
		UserName: "aaa",
	}, "12345", 8)
	fmt.Println(token, err)
}

func TestParseToken(t *testing.T) {
	payload, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJhYWEiLCJyb2xlIjoxLCJleHAiOjE3MzE4Mzk3NjZ9.NFDFr6ws6ap2lcsmKJAlleaUy-qMj6uaKXZqJyvsyqo", "12345")
	fmt.Println(payload, err)
}
