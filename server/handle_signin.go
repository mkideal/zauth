package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleSignin(w http.ResponseWriter, r *http.Request) {
	argv := new(api.SigninReq)
	err := argv.Parse(r)
	if err != nil {
		log.Warn("Signin parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signin request, IP=%v", httputil.IP(r))
	account := model.JoinAccount(model.AccountType(argv.AccountType), argv.Account)
	if account == "" {
		svr.response(w, http.StatusBadRequest, "invalid accountType or account")
	}
	user, err := svr.userRepo.GetUserByAccount(account)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		svr.responseErrorCode(w, api.ErrorCode_UserNotFound, "user-not-found")
		return
	}
	if !model.ValidatePassword(user, argv.Password) {
		svr.responseErrorCode(w, api.ErrorCode_IncorrectPassword, "incorrect-password")
		return
	}
	_, err = svr.setSession(w, r, user.Id)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	svr.response(w, http.StatusOK, makeUserInfo(user))
}
