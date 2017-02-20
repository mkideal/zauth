package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/Bowery/prompt"
	"github.com/google/shlex"
	"github.com/mkideal/cli"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/api/golang"
)

type Context struct {
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
	UserRes           api.UserRes
}

func (ctx *Context) outputJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(v)
}

func (ctx *Context) onError(err error) bool {
	e := golang.ErrorResponse(err)
	if e != nil {
		ctx.outputJSON(e)
	} else if err != nil {
		fmt.Printf("Warn: %v", err)
	}
	return err != nil
}

func (ctx *Context) onHelpRes(res *api.HelpRes, err error) {
	if ctx.onError(err) {
		return
	}
	ctx.outputJSON(res)
}

var context = &Context{}

func printFields(v reflect.Value, owner string) {
	return
	t := v.Type()
	fmt.Printf("numField: %d\n", t.NumField())
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("%dth field of %s: %s\n", i, owner, t.Field(i).Name)
	}
}

// $(xxx)
func value(s string) (res interface{}) {
	defer func() {
		if e := recover(); e != nil {
			res = e
		}
	}()

	fields := strings.Split(s, ".")
	v := reflect.ValueOf(context)
	parent := "Context"
	for i, field := range fields {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if i > 0 {
			printFields(v, parent)
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
	golang.Config
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		client := golang.NewClient(argv.Config)
		_ = client
		quit := false
		for !quit {
			line, err := prompt.Basic(prefix(context), false)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				break
			}
			line = strings.TrimSpace(line)
			args, err := shlex.Split(line)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				break
			}
			if len(args) == 0 || args[0] == "" {
				continue
			}
			cmd := strings.ToLower(args[0])
			switch cmd {
			case "exit", "quit", "q":
				fmt.Printf("bye~\n")
				quit = true
			case "exist", "account_exist":
				//client.AccountExist(req)
			case "auto_signup":
			case "signup":
			case "signin":
			case "signout":
			case "token":
			case "token_auth":
			case "user":
			case "help":
				req := new(api.HelpReq)
				context.onHelpRes(client.Help(req))
			default:
				fmt.Printf("Unknown command `%s`\n", cmd)
			}
		}
		return nil
	})
}
