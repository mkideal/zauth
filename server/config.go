package server

import (
	"errors"
	"github.com/mkideal/cli"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Driver                     string `cli:"driver" usage:"sql database driver: mysql" dft:"mysql"`
	DataSourceName             string `cli:"dsn" usage:"data source name for specified driver" dft:"$ACCOUNT_DSN"`
	ThirdParty                 string `cli:"third-party" usage:"third party modules which seperated by ,"`
	Addr                       string `cli:"addr" usage:"HTTP address" dft:"127.0.0.1:5200"`
	Mode                       string `cli:"m,mode" usage:"running mode: debug/release" dft:"release"`
	CookieKey                  string `cli:"cookie" usage:"cookie key" dft:"authd"`
	SessionExpireDuration      int64  `cli:"session-expire-duration" usage:"session expire duration(seconds)" dft:"3600"`
	HTMLDir                    string `cli:"html" usage:"HTML static directory" dft:"html"`
	HTMLRoouter                string `cli:"html-router" usage:"HTML static files router" dft:"/"`
	TelnoVerifyCodeMaxInterval int64  `cli:"sms-max-interval" usage:"SMS code max interval seconds" dft:"60"`
	TelnoVerifyCodeExpiration  int64  `cli:"sms-expiration" usage:"SMS code expiration seconds" dft:"300"`
	TelnoVerifyCodeLength      int    `cli:"sms-length" usage:"length of SMS code, must be in range [3,8]" dft:"6"`
	EnableTelnoVerify          bool   `cli:"sms-enable" usage:"enable telno verify" dft:"true"`
	SendTelnoVerifyCode        bool   `cli:"sms-send" usage:"send telno code verify message or not" dft:"true"`

	SMS `cli:"-"`

	Pages
	WhiteTelnoList []string `cli:"white-telno" usage:"white telno list"`
}

func (config Config) Validate(ctx *cli.Context) error {
	if config.TelnoVerifyCodeLength < 3 || config.TelnoVerifyCodeLength > 8 {
		return errors.New("SMS code length must be range in [3,8]")
	}
	return nil
}

func (config Config) IsWhiteTelno(telno string) bool {
	if !config.EnableTelnoVerify {
		return true
	}
	for _, x := range config.WhiteTelnoList {
		if x == telno {
			return true
		}
	}
	return false
}

type SMS struct {
	SMSURL       string
	SMSUsername  string
	SMSPassword  string
	SMSMsgFormat string
}

type Pages struct {
	Authorize string `cli:"page-authorize" usage:"web page URL for authorize" dft:"/authorize.html"`
	Login     string `cli:"page-login" usage:"web page URL for login" dft:"/login.html"`
}
