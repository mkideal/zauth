package server

import (
	"net/http"
	"time"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
)

func (svr *Server) handleTokenAuth(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.TokenAuthReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("TokenAuth parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("TokenAuth request, IP=%v", ip)

	if accessToken := svr.getTokenFromHeader(r); accessToken != "" {
		argv.AccessToken = accessToken
	}

	token, err := svr.tokenRepo.GetToken(argv.AccessToken)
	if err != nil {
		log.Error("%s: get token %s error: %v", argv.CommandName(), argv.AccessToken, err)
		svr.errorResponse(w, r, err)
		return
	}
	if token == nil {
		log.Info("%s: token %s not found", argv.CommandName(), argv.AccessToken)
		svr.errorResponse(w, r, api.ErrorCode_TokenNotFound)
		return
	}
	user, err := svr.userRepo.GetUser(token.Uid)
	if err != nil {
		log.Error("%s: get user %d error: %v", argv.CommandName(), token.Uid, err)
		svr.errorResponse(w, r, err)
		return
	}
	if user == nil {
		log.Warn("%s: user %d not found", argv.CommandName(), token.Uid)
		svr.errorResponse(w, r, api.ErrorCode_UserNotFound)
		return
	}
	svr.response(w, r, api.TokenAuthRes{
		User:  makeUserInfo(user),
		Token: makeTokenInfo(token),
	})
	user.LastLoginAt = model.FormatTime(time.Now())
	user.LastLoginIp = ip
	meta := user.Meta()
	svr.userRepo.UpdateUser(user, meta.F_last_login_at, meta.F_last_login_ip)
}
