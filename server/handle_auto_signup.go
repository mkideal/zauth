package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleAutoSignup(w http.ResponseWriter, r *http.Request) {
	argv := new(api.AutoSignupReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AutoSignup parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AutoSignup request, IP=%v", httputil.IP(r))
	client := svr.clientAuth(argv.CommandName(), w, r)
	if client == nil {
		// NOTE: response returned in clientAuth
		return
	}
	user := new(model.User)
	user.AccountType = model.AccountType_Auto
	user.CreatedIp = httputil.IP(r)
	if err := svr.userRepo.AddUser(user, ""); err != nil {
		log.Error("%s: add user error: %v", argv.CommandName(), err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	accessToken, err := svr.tokenRepo.NewToken(client, user, "")
	if err != nil {
		log.Error("%s: new token error: %v", argv.CommandName(), err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	res := api.AutoSignupRes{
		Uid:          user.Id,
		AccessToken:  accessToken.Token,
		RefreshToken: accessToken.RefreshToken,
		ExpireAt:     accessToken.ExpireAt,
	}
	svr.response(w, http.StatusOK, res)
}
