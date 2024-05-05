package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const (
	// InternalErrorType type of common errors
	InternalErrorType = "INTERNAL"
	// GenericErrorType type of common errors
	GenericErrorType = "GENERIC"
	// ValidationErrorType type of common errors
	ValidationErrorType = "VALIDATION"
)

// ErrorResponse represents the error response
type ErrorResponse struct {
	Error *HTTPError `json:"error"`
}

// HTTPError represents an error that occurred while handling a request
type HTTPError struct {
	Code     int             `json:"code"`
	Type     string          `json:"type"`
	Message  string          `json:"message"`
	Internal error           `json:"-"`
	Extra    *HTTPErrorExtra `json:"extra,omitempty"`
}

// HTTPErrorExtra represents an error that occurred while handling a request with extra information
type HTTPErrorExtra struct {
	ButtonLabel string `json:"button_label"`
	Title       string `json:"title"`
	Action      string `json:"action"`
	IconType    string `json:"icon_type"`
}

// NewHTTPError creates a new HTTPError instance
func NewHTTPError(code int, etype string, message ...string) *HTTPError {
	// * message[0] - message
	// * message[1] - title
	// * message[2] - label of button
	// * message[3] - action of button
	// * message[4] - icon of button

	he := new(HTTPError)
	he.Code = code
	he.Type = etype

	if len(message) > 0 {
		he.Message = message[0]
		switch len(message) {
		case 2:
			he.Extra = &HTTPErrorExtra{
				Title: message[1],
			}
		case 3:
			he.Extra = &HTTPErrorExtra{
				Title:       message[1],
				ButtonLabel: message[2],
			}
		case 4:
			he.Extra = &HTTPErrorExtra{
				Title:       message[1],
				ButtonLabel: message[2],
				Action:      message[3],
			}
		case 5:
			he.Extra = &HTTPErrorExtra{
				Title:       message[1],
				ButtonLabel: message[2],
				Action:      message[3],
				IconType:    message[4],
			}
		}

	} else {
		he.Message = http.StatusText(code)
	}
	return he
}

// NewHTTPInternalError creates a new HTTPError instance for internal error
func NewHTTPInternalError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusInternalServerError, Type: InternalErrorType, Message: message}
}

// NewHTTPGenericError creates a new HTTPError instance for generic error
func NewHTTPGenericError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Type: GenericErrorType, Message: message}
}

// NewHTTPValidationError creates a new HTTPError instance for validation error
func NewHTTPValidationError(message string) *HTTPError {
	return &HTTPError{Code: http.StatusBadRequest, Type: ValidationErrorType, Message: message}
}

// Error makes it compatible with `error` interface
func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, type=%s, message=%s", he.Code, he.Type, he.Message)
}

// SetInternal sets actual internal error for more details
func (he *HTTPError) SetInternal(err error) *HTTPError {
	he.Internal = err
	return he
}

// SetType sets custom type error for more details
func (he *HTTPError) SetType(errType string) *HTTPError {
	he.Type = errType
	return he
}

// SetTitle sets custom title error for more details
func (he *HTTPError) SetTitle(title string) *HTTPError {
	he.Extra = &HTTPErrorExtra{Title: title}
	return he
}

// SetExtraData sets custom extra data error for button when error
func (he *HTTPError) SetExtraData(title, buttonLabel, action, iconType string) *HTTPError {
	he.Extra = &HTTPErrorExtra{Title: title, ButtonLabel: buttonLabel, Action: action, IconType: iconType}
	return he
}

// SetButtonLabel sets custom button label internal error for more details
func (he *HTTPError) SetButtonLabel(label string) *HTTPError {
	he.Extra = &HTTPErrorExtra{ButtonLabel: label}
	return he
}

// SetMessage sets custom message error for more details
func (he *HTTPError) SetMessage(message string) *HTTPError {
	he.Message = message
	return he
}

// ErrorHandler represents the custom http error handler
type ErrorHandler struct {
	e *echo.Echo
}

// NewErrorHandler returns the ErrorHandler instance
func NewErrorHandler(e *echo.Echo) *ErrorHandler {
	return &ErrorHandler{e}
}

// Handle is a centralized HTTP error handler.
func (ce *ErrorHandler) Handle(err error, c echo.Context) {
	httpErr := NewHTTPError(http.StatusInternalServerError, InternalErrorType)

	switch e := err.(type) {
	case *HTTPError:
		if e.Code != 0 {
			httpErr.Code = e.Code
		}
		if e.Type != "" {
			httpErr.Type = e.Type
		} else {
			httpErr.Type = GenericErrorType
		}
		if e.Message != "" {
			httpErr.Message = e.Message
		}
		if e.Internal != nil && !c.Response().Committed {
			fmt.Printf("internal err: %+v", e.Internal)
		}
		if e.Extra != nil {
			httpErr.Extra = e.Extra
		}

	case *echo.HTTPError:
		httpErr.Code = e.Code
		httpErr.Type = GenericErrorType
		switch em := e.Message.(type) {
		case string:
			httpErr.Message = em
		case []string:
			httpErr.Message = strings.Join(em, "\n")
		case map[string]interface{}:
			if jsonStr, err := json.Marshal(em); err == nil {
				httpErr.Message = string(jsonStr)
			}
		default:
			httpErr.Message = fmt.Sprintf("%+v", em)
		}
		if e.Internal != nil && !c.Response().Committed {
			fmt.Printf("internal err: %+v", e.Internal)
		}

	case validator.ValidationErrors:
		httpErr.Code = http.StatusBadRequest
		httpErr.Type = ValidationErrorType
		var errMsg []string
		for _, v := range e {
			errMsg = append(errMsg, getVldErrorMsg(v))
		}
		httpErr.Message = strings.Join(errMsg, "\n")
	default:
		if ce.e.Debug {
			httpErr.Message = err.Error()
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(httpErr.Code)
		} else {
			err = c.JSON(httpErr.Code, ErrorResponse{Error: httpErr})
		}
		if err != nil {
			fmt.Printf("internal err: %+v", err)
		}
	}
}

var validationErrors = map[string]string{
	"required": " is required, but was not received",
	"min":      "'s value or length is less than allowed",
	"max":      "'s value or length is bigger than allowed",
	"date":     "'s value should be in form of YYYY-MM-DD",
	"email":    "'s value should be a valid email address",
	"mobile":   "'s value should be a valid mobile number",
	"url":      "'s value should be a valid URL",
}

// CustomValidator holds custom validator
type CustomValidator struct {
	V *validator.Validate
}

// NewValidator creates new custom validator
func NewValidator() *CustomValidator {
	V := validator.New()
	return &CustomValidator{V}
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.V.Struct(i)
}

func getVldErrorMsg(v validator.FieldError) string {
	field := v.Field()
	vtag := v.ActualTag()
	vtagVal := v.Param()

	if msg, ok := validationErrors[vtag]; ok {
		return field + msg
	}

	switch vtag {
	case "oneof":
		return field + " should be one of " + strings.Replace(vtagVal, " ", ", ", -1)
	case "ltfield":
		return field + " should be less than " + vtagVal
	case "gtfield":
		return field + " should be greater than " + vtagVal
	case "eqfield":
		return field + " does not match " + vtagVal
	}

	return field + " failed on " + vtag + " validation"
}
