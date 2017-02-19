package oauth2

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	TokenType         = "bearer"
	TokenHeaderPrefix = "Bearer "
)

// Grant type
const (
	GrantAuthenticationCode = "authorization_code"
	GrantPassword           = "password"
	GrantRefreshToken       = "refresh_token"
)

//	Response type
const (
	ResponseCode        = "code"
	ResponseAccessToken = "access_token"
)

type Error struct {
	statusCode  int    `json:"-"`
	Code        string `json:"error"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`
	State       string `json:"state,omitempty"`
}

func (e Error) Status() int {
	if e.statusCode == 0 {
		return http.StatusBadRequest
	}
	return e.statusCode
}

func (e *Error) SetStatus(statusCode int) {
	e.statusCode = statusCode
}

func (e Error) Encode() string {
	return e.EncodeWith(nil)
}

func (e Error) EncodeWith(values map[string]interface{}) string {
	others := ""
	if values != nil {
		for k, v := range values {
			others += fmt.Sprintf("%s=%s&", k, template.HTMLEscapeString(fmt.Sprintf("%v", v)))
		}
	}
	return fmt.Sprintf("%serror=%s&error_description=%s&error_uri=%s&state=%s",
		others,
		e.Code,
		template.HTMLEscapeString(e.Description),
		template.HTMLEscapeString(e.URI),
		template.HTMLEscapeString(e.State),
	)
}

func (e Error) Error() string {
	if e.Description == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

// Error code
type OAuthErrorCode string

func (e OAuthErrorCode) Error() string   { return string(e) }
func (e OAuthErrorCode) NewError() Error { return NewError(string(e), "") }

const (
	ErrorInvalidRequest          OAuthErrorCode = "invalid_request"
	ErrorInvalidClient           OAuthErrorCode = "invalid_client"
	ErrorInvalidScope            OAuthErrorCode = "invalid_scope"
	ErrorInvalidGrant            OAuthErrorCode = "invalid_grant"
	ErrorUnauthorizedClient      OAuthErrorCode = "unauthorized_client"
	ErrorAccessDenied            OAuthErrorCode = "access_denied"
	ErrorUnsupportedResponseType OAuthErrorCode = "unsupported_response_type"
	ErrorServerError             OAuthErrorCode = "server_error"
	ErrorTemporarilyUnavailable  OAuthErrorCode = "temporarily_unavailable"
	ErrorUnsupportedGrantType    OAuthErrorCode = "unsupported_grant_type"
)

func NewError(code, description string) Error {
	return Error{
		Code:        code,
		Description: description,
	}
}

func WrapError(err error) Error {
	if authErr, ok := err.(Error); ok {
		return authErr
	}
	e := NewError(string(ErrorServerError), err.Error())
	e.statusCode = http.StatusInternalServerError
	return e
}
