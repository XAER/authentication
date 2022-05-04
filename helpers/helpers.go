package helpers

import (
	"net/mail"
	"os"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IsEmailOrUsername(identifier string) string {
	_, err := mail.ParseAddress(identifier)
	if err != nil {
		return "username"
	}
	return "email"
}
