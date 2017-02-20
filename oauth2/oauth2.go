package oauth2

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

// Error code
type OAuthErrorCode string

func (e OAuthErrorCode) Error() string { return string(e) }

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
