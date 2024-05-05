package http

import (
	"music-master/internal/util/server"

	"github.com/labstack/echo/v4"
)

// ReqID returns id url parameter.
func ReqID(c echo.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", server.NewHTTPValidationError("Invalid ID")
	}
	return id, nil
}

// ListRequest holds data of listing request from react-admin
// swagger:model ListRequest
type ListRequest struct {
	// Number of records per page
	// default: 25
	Limit int `json:"l,omitempty" query:"l" validate:"max=300"`
	// Current page number
	// default: 1
	Page int `json:"p,omitempty" query:"p"`
	// JSON string of filter. E.g: {"field_name":"value"}
	// default:
	Filter string `json:"f,omitempty" query:"f"`
}
