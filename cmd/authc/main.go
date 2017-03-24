package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/Bowery/prompt"
	"github.com/google/shlex"
	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/typeconv"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/api/go/authc"
)

type Context struct {
	Addr  string
	Token api.TokenInfo
	User  api.UserInfo
	Error error

	// temporary data

	TokenRes          api.TokenRes
	AuthorizeCheckRes api.AuthorizeCheckRes
	TokenAuthRes      api.TokenAuthRes
	SignupRes         api.SignupRes
	AutoSignupRes     api.AutoSignupRes
	AccountExistRes   api.AccountExistRes
	SigninRes         api.SigninRes
	SignoutRes        api.SignoutRes
	SMSCodeRes        api.SMSCodeRes
	TwoFactorAuthRes  api.TwoFactorAuthRes
	UserRes           api.UserRes
}

func (ctx *Context) outputJSON(v interface{}, title ...string) {
	for _, t := range title {
		ctx.println(t)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(v)
}

func (ctx *Context) println(s string) {
	fmt.Println(s)
}

func (ctx *Context) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (ctx *Context) parseRequest(args []string, req interface{}) error {
	return cli.Parse(args, req)
}

func (ctx *Context) onError(err error) bool {
	e := authc.ErrorResponse(err)
	if e != nil {
		ctx.outputJSON(e)
	} else if err != nil {
		ctx.printf("Warn: %v", err)
	}
	return err != nil
}

func (ctx *Context) onHelpRes(res *api.HelpRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "Help response")
}

func (ctx *Context) onAccountExist(res *api.AccountExistRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "AccountExist response")
	ctx.AccountExistRes = *res
}

func (ctx *Context) onAutoSignup(res *api.AutoSignupRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "AutoSignup response")
	ctx.AutoSignupRes = *res
	ctx.Token = res.Token
	ctx.User.Id = res.Uid
}

func (ctx *Context) onSignup(res *api.SignupRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "Signup response")
	ctx.SignupRes = *res
	ctx.User = res.User
	ctx.Token = res.Token
}

func (ctx *Context) onSignin(res *api.SigninRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "Signin response")
	ctx.SigninRes = *res
	ctx.Token = res.Token
	ctx.User = res.User
}

func (ctx *Context) onSignout(res *api.SignoutRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "Signout response")
	ctx.SignoutRes = *res
	ctx.User = api.UserInfo{}
	ctx.Token = api.TokenInfo{}
}

func (ctx *Context) onToken(res *api.TokenRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "Token response")
	ctx.TokenRes = *res
	ctx.Token = res.Token
}

func (ctx *Context) onTokenAuth(res *api.TokenAuthRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "TokenAuth response")
	ctx.TokenAuthRes = *res
	ctx.Token = res.Token
	ctx.User = res.User
}

func (ctx *Context) onSMSCode(res *api.SMSCodeRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "SMSCode response")
	ctx.SMSCodeRes = *res
}

func (ctx *Context) onTwoFactorAuth(res *api.TwoFactorAuthRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "TwoFactorAuth response")
	ctx.TwoFactorAuthRes = *res
	ctx.User = res.User
	ctx.Token = res.Token
}

func (ctx *Context) onUser(res *api.UserRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res, "User response")
	ctx.UserRes = *res
	ctx.User = res.User
}

var context = &Context{}

// $xxx
func value(s string) (res interface{}) {
	defer func() {
		if e := recover(); e != nil {
			res = e
		}
	}()

	fields := strings.Split(s, ".")
	v := reflect.ValueOf(context)
	parent := "Context"
	for _, field := range fields {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		next := v.FieldByName(field)
		if !next.IsValid() {
			field = strings.Title(field)
		}
		next = v.FieldByName(field)
		if !next.IsValid() {
			panic(fmt.Sprintf("field %s not found in %s", field, parent))
		}
		parent += "." + field
		v = next
	}
	res = v.Interface()
	return
}

func prefix(ctx *Context) string {
	return fmt.Sprintf("~%d$ ", ctx.User.Id)
}

type argT struct {
	cli.Helper
	authc.Config
	Address string `cli:"addr" usage:"authd http address"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		client := authc.NewClient(argv.Config)
		context.Addr = argv.Address
		quit := false
		for !quit {
			line, err := prompt.Basic(prefix(context), false)
			if err != nil {
				context.onError(err)
				break
			}
			err, quit = execLine(client, line)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		return nil
	})
}

func execLine(client *authc.Client, line string) (err error, quit bool) {
	line = strings.TrimSpace(line)
	args, err := shlex.Split(line)
	if err != nil {
		context.printf("Error: %v\n", err)
		quit = true
		return
	}
	if len(args) == 0 || args[0] == "" {
		return
	}
	cmd := strings.ToLower(args[0])
	args = args[1:]
	for i, arg := range args {
		if strings.HasPrefix(arg, "$") {
			args[i] = typeconv.ToString(value(strings.TrimPrefix(arg, "$")))
		}
	}
	switch cmd {
	case "exit", "quit", "q":
		context.printf("bye~\n")
		quit = true
	case "exist", "account_exist":
		req := new(api.AccountExistReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onAccountExist(client.AccountExist(context.Addr, req))
		}
	case "auto_signup":
		req := new(api.AutoSignupReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onAutoSignup(client.AutoSignup(context.Addr, req))
		}
	case "signup":
		req := new(api.SignupReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onSignup(client.Signup(context.Addr, req))
		}
	case "signin":
		req := new(api.SigninReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onSignin(client.Signin(context.Addr, req))
		}
	case "signout":
		req := new(api.SignoutReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onSignout(client.Signout(context.Addr, req))
		}
	case "token":
		req := new(api.TokenReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onToken(client.Token(context.Addr, req))
		}
	case "token_auth":
		req := new(api.TokenAuthReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onTokenAuth(client.TokenAuth(context.Addr, req))
		}
	case "sms", "smscode":
		req := new(api.SMSCodeReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onSMSCode(client.SMSCode(context.Addr, req))
		}
	case "2fa", "2fa_auth":
		req := new(api.TwoFactorAuthReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onTwoFactorAuth(client.TwoFactorAuth(context.Addr, req))
		}
	case "user":
		req := new(api.UserReq)
		err = context.parseRequest(args, req)
		if err == nil {
			context.onUser(client.User(context.Addr, req))
		}
	case "help":
		req := new(api.HelpReq)
		context.onHelpRes(client.Help(context.Addr, req))
	case "p", "print":
		for _, arg := range args {
			context.outputJSON(value(arg))
		}
	case "exec":
		for _, filename := range args {
			var (
				content []byte
				index   int
			)
			content, err = ioutil.ReadFile(filename)
			if err != nil {
				return
			}
			for {
				var (
					token   []byte
					advance int
				)
				advance, token, err = bufio.ScanLines(content[index:], false)
				if advance == 0 {
					break
				}
				index += advance
				if err != nil {
					return
				}
				err, quit = execLine(client, string(token))
				if quit {
					return
				}
				if err != nil {
					return
				}
			}
		}
	default:
		context.printf("Unknown command `%s`\n", cmd)
	}
	return
}
