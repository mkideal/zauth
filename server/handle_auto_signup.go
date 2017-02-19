package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleAutoSignup(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.AutoSignupReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AutoSignup parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("AutoSignup request, IP=%v", ip)
	client := svr.clientAuth(argv.CommandName(), w, r)
	if client == nil {
		// NOTE: response returned in clientAuth
		return
	}
	user := new(model.User)
	user.AccountType = model.AccountType_Auto
	user.CreatedIp = ip
	if err := svr.userRepo.AddUser(user, ""); err != nil {
		log.Error("%s: add user error: %v", argv.CommandName(), err)
		svr.errorResponse(w, err)
		return
	}
	token, err := svr.tokenRepo.NewToken(user, client.Id, client.Scope)
	if err != nil {
		log.Error("%s: new token error: %v", argv.CommandName(), err)
		svr.errorResponse(w, err)
		return
	}
	res := api.AutoSignupRes{
		Uid:   user.Id,
		Token: makeTokenInfo(token),
	}
	svr.response(w, res)
}
