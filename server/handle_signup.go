package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	argv := new(api.SignupReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Signup parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signup request, IP=%v", httputil.IP(r))
	user := new(model.User)
	user.AccountType = model.AccountType(argv.AccountType)
	// TODO: 检查 accountType,暂时不支持第三方账号注册
	user.CreatedIP = httputil.IP(r)
	user.Account = argv.Account
	if err := svr.userRepo.AddUser(user, argv.Password); err != nil {
		log.Error("%s: add user %s error: %v", argv.CommandName(), user.Account, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	res := api.SignupRes{
		Uid:      user.Id,
		Account:  user.Account,
		Nickname: user.Nickname,
	}
	svr.response(w, http.StatusOK, res)
}
