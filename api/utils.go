package api

import (
	"fmt"
	"html/template"
	"net/http"

	"bitbucket.org/mkideal/accountd/oauth2"
)

const (
	TwoFaType_Telno = "telno"
	TwoFaType_Email = "email"
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
	e := NewError(string(oauth2.ErrorServerError), err.Error())
	e.statusCode = http.StatusInternalServerError
	return e
}
