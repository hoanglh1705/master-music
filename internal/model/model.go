package model

import "github.com/labstack/echo"

// AuthUser represents data stored in JWT token for customer
type AuthUser struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

// Auth represents auth interface
type Auth interface {
	User(echo.Context) *AuthUser
}
