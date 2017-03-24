package authc

import (
	"encoding/json"
	"net/http"
	"strings"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/oauth2"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Config struct {
	ClientId     string `cli:"client-id" usage:"client id"`
	ClientSecret string `cli:"client-secret" usage:"client secret"`
	CookieName   string `cli:"cookie" usage:"cookie name" dft:"accountd"`
	Router

	HTTPClient HTTPClient `cli:"-" json:"-"`
}

type Router struct {
	AccountExist  string `cli:"r-account-exist" usage:"router of command AccountExist"`
	Authorize     string `cli:"r-authorize" usage:"router of command Authorize"`
	AutoSignup    string `cli:"r-auto-signup" usage:"router of command AutoSignup"`
	Signup        string `cli:"r-signup" usage:"router of command Signup"`
	Signin        string `cli:"r-signin" usage:"router of command Signin"`
	Signout       string `cli:"r-signout" usage:"router of command Signout"`
	Token         string `cli:"r-token" usage:"router of command Token"`
	TokenAuth     string `cli:"r-token-auth" usage:"router of command TokenAuth"`
	SMSCode       string `cli:"r-smscode" usage:"router of command SMSCode"`
	TwoFactorAuth string `cli:"r-2fa-auth" usage:"router of command TwoFactorAuth"`
	User          string `cli:"r-user" usage:"router of command User"`
	Help          string `cli:"r-help" usage:"router of command Help"`
}

func (r *Router) Init() {
	initString(&r.AccountExist, "/v1/account_exist")
	initString(&r.Authorize, "/v1/authorize")
	initString(&r.AutoSignup, "/v1/auto_signup")
	initString(&r.Signup, "/v1/signup")
	initString(&r.Signin, "/v1/signin")
	initString(&r.Signout, "/v1/signout")
	initString(&r.Token, "/v1/token")
	initString(&r.TokenAuth, "/v1/token_auth")
	initString(&r.SMSCode, "/v1/smscode")
	initString(&r.TwoFactorAuth, "/v1/two_factor_auth")
	initString(&r.User, "/v1/user")
	initString(&r.Help, "/v1/help")
}

func initString(ptr *string, dft string) {
	if *ptr == "" {
		*ptr = dft
	}
}

func ErrorResponse(err error) *api.Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(api.Error); ok {
		return &e
	}
	if e, ok := err.(*api.Error); ok {
		return e
	}
	return nil
}

type Client struct {
	config Config

	StatusCode int
	Cookie     *http.Cookie
}

const (
	tagPrefix          = "__tag_"
	tagClientBasicAuth = tagPrefix + "client_basic_auth"
)

func NewClient(config Config) *Client {
	c := &Client{
		config: config,
	}
	if c.config.HTTPClient == nil {
		c.config.HTTPClient = http.DefaultClient
	}
	c.config.Router.Init()
	return c
}

func (client *Client) url(address, router string) string {
	if !strings.HasPrefix(router, "/") {
		return router
	}
	if strings.HasSuffix(address, "/") {
		return address + strings.TrimPrefix(router, "/")
	}
	return address + router
}

func (client *Client) beforeDoHTTP(r *http.Request, headers map[string]string) {
	if client.Cookie != nil {
		r.AddCookie(client.Cookie)
	}
	if _, found := headers[tagClientBasicAuth]; found {
		r.SetBasicAuth(client.config.ClientId, client.config.ClientSecret)
	}
	if headers != nil {
		for k, v := range headers {
			if !strings.HasPrefix(k, tagPrefix) {
				r.Header.Set(k, v)
			}
		}
	}
}

func (client *Client) get(url string, req api.Request, res interface{}, headers map[string]string) error {
	url += "?" + req.Values().Encode()
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client.beforeDoHTTP(r, headers)
	resp, err := client.config.HTTPClient.Do(r)
	return client.handleResponse(resp, err, res)
}

func (client *Client) post(url string, req api.Request, res interface{}, headers map[string]string) error {
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(req.Values().Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client.beforeDoHTTP(r, headers)
	resp, err := client.config.HTTPClient.Do(r)
	return client.handleResponse(resp, err, res)
}

func (client *Client) handleResponse(resp *http.Response, err error, res interface{}) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == client.config.CookieName {
			client.Cookie = cookie
			break
		}
	}
	client.StatusCode = resp.StatusCode
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err := api.Error{}
		if e := decoder.Decode(&err); e != nil {
			return e
		}
		return err
	}
	return decoder.Decode(res)
}

func (client *Client) AccountExist(addr string, req *api.AccountExistReq) (res *api.AccountExistRes, err error) {
	res = new(api.AccountExistRes)
	err = client.get(client.url(addr, client.config.Router.AccountExist), req, res, nil)
	return
}

func (client *Client) AutoSignup(addr string, req *api.AutoSignupReq) (res *api.AutoSignupRes, err error) {
	res = new(api.AutoSignupRes)
	err = client.post(client.url(addr, client.config.Router.AutoSignup), req, res, map[string]string{
		tagClientBasicAuth: "true",
	})
	return
}

func (client *Client) Signup(addr string, req *api.SignupReq) (res *api.SignupRes, err error) {
	res = new(api.SignupRes)
	err = client.post(client.url(addr, client.config.Router.Signup), req, res, nil)
	return
}

func (client *Client) Signin(addr string, req *api.SigninReq) (res *api.SigninRes, err error) {
	res = new(api.SigninRes)
	err = client.post(client.url(addr, client.config.Router.Signin), req, res, nil)
	return
}

func (client *Client) Signout(addr string, req *api.SignoutReq) (res *api.SignoutRes, err error) {
	res = new(api.SignoutRes)
	err = client.post(client.url(addr, client.config.Router.Signout), req, res, nil)
	return
}

func (client *Client) Token(addr string, req *api.TokenReq) (res *api.TokenRes, err error) {
	res = new(api.TokenRes)
	err = client.post(client.url(addr, client.config.Router.Token), req, res, map[string]string{
		tagClientBasicAuth: "true",
	})
	return
}

func (client *Client) TokenAuth(addr string, req *api.TokenAuthReq) (res *api.TokenAuthRes, err error) {
	res = new(api.TokenAuthRes)
	err = client.post(client.url(addr, client.config.Router.TokenAuth), req, res, map[string]string{
		"Authorization": oauth2.TokenHeaderPrefix + req.AccessToken,
	})
	return
}

func (client *Client) SMSCode(addr string, req *api.SMSCodeReq) (res *api.SMSCodeRes, err error) {
	res = new(api.SMSCodeRes)
	err = client.post(client.url(addr, client.config.Router.SMSCode), req, res, nil)
	return
}

func (client *Client) TwoFactorAuth(addr string, req *api.TwoFactorAuthReq) (res *api.TwoFactorAuthRes, err error) {
	res = new(api.TwoFactorAuthRes)
	err = client.post(client.url(addr, client.config.Router.TwoFactorAuth), req, res, nil)
	return
}

func (client *Client) User(addr string, req *api.UserReq) (res *api.UserRes, err error) {
	res = new(api.UserRes)
	err = client.get(client.url(addr, client.config.Router.User), req, res, nil)
	return
}

func (client *Client) Help(addr string, req *api.HelpReq) (res *api.HelpRes, err error) {
	res = new(api.HelpRes)
	err = client.get(client.url(addr, client.config.Router.Help), req, res, nil)
	return
}
