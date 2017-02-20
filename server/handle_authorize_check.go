package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleAuthorizeCheck(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.AuthorizeCheckReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AuthorizeCheck parse arguments error: %v, IP=%v", ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("AuthorizeCheck request, IP=%v", ip)
	if argv.ClientId == "" {
		log.Info("%s: missing argument client_id", argv.CommandName())
		svr.errorResponse(w, r, api.ErrorCode_MissingArgument.NewError("missing client_id"))
		return
	}
	client, err := svr.clientRepo.GetClient(argv.ClientId)
	if err != nil {
		log.Error("%s: get client %s error: %v", argv.CommandName(), argv.ClientId, err)
		svr.errorResponse(w, r, err)
		return
	}
	if client == nil {
		log.Info("%s: client %s not found", argv.CommandName(), argv.ClientId)
		svr.errorResponse(w, r, api.ErrorCode_ClientNotFound)
		return
	}
	session := svr.getSession(r)
	if session == nil {
		log.Info("%s: session not found", argv.CommandName())
		svr.errorResponse(w, r, api.ErrorCode_SessionNotFound)
		return
	}
	user, err := svr.userRepo.GetUser(session.Uid)
	if err != nil {
		log.Error("%s: get user %d error: %v", argv.CommandName(), session.Uid, err)
		svr.errorResponse(w, r, err)
		return
	}
	if user == nil {
		log.Warn("%s: user %d not found", argv.CommandName(), session.Uid)
		svr.errorResponse(w, r, api.ErrorCode_UserNotFound)
		return
	}
	svr.response(w, r, api.AuthorizeCheckRes{
		Application: client.Name,
		Username:    user.Nickname,
	})
}
