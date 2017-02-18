package oauth2

import (
	"fmt"
	"html/template"
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
	Code        string `json:"error"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`
	State       string `json:"state,omitempty"`
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
const (
	ErrorInvalidRequest          = "invalid_request"
	ErrorInvalidClient           = "invalid_client"
	ErrorInvalidScope            = "invalid_scope"
	ErrorInvalidGrant            = "invalid_grant"
	ErrorUnauthorizedClient      = "unauthorized_client"
	ErrorAccessDenied            = "access_denied"
	ErrorUnsupportedResponseType = "unsupported_response_type"
	ErrorServerError             = "server_error"
	ErrorTemporarilyUnavailable  = "temporarily_unavailable"
	ErrorUnsupportedGrantType    = "unsupported_grant_type"
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
	return NewError(ErrorServerError, err.Error())
}
