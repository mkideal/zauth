package server

import (
	"net/http"
	"time"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleSignin(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.SigninReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Signin parse arguments error: %v, IP=%v", err, ip)
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signin request, IP=%v", ip)
	account := model.JoinAccount(model.AccountType(argv.AccountType), argv.Account)
	if account == "" {
		log.Info("%s: missing account_type or account", argv.CommandName())
		svr.response(w, http.StatusBadRequest, "invalid accountType or account")
		return
	}
	user, err := svr.userRepo.GetUserByAccount(account)
	if err != nil {
		log.Error("%s: get user by account %s error: %v", argv.CommandName(), account, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		log.Info("%s: account %s not found", argv.CommandName(), account)
		svr.responseErrorCode(w, api.ErrorCode_UserNotFound, "user-not-found")
		return
	}
	if !model.ValidatePassword(user, argv.Password) {
		log.Info("%s: incorrect password for user (%d,%s)", argv.CommandName(), user.Id, user.Account)
		svr.responseErrorCode(w, api.ErrorCode_IncorrectPassword, "incorrect-password")
		return
	}
	_, err = svr.setSession(w, r, user.Id)
	if err != nil {
		log.Error("%s: set session error: %v", argv.CommandName(), err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	token, err := svr.tokenRepo.NewToken(user, "", "")
	if err != nil {
		log.Error("%s: new token error: %v", argv.CommandName(), err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	svr.response(w, http.StatusOK, api.SigninRes{
		User:  makeUserInfo(user),
		Token: makeTokenInfo(token),
	})
	user.LastLoginAt = model.FormatTime(time.Now())
	user.LastLoginIp = ip
	meta := model.UserMetaVar
	svr.userRepo.UpdateUser(user, meta.F_last_login_at, meta.F_last_login_ip)
}
