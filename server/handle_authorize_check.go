package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleAuthorizeCheck(w http.ResponseWriter, r *http.Request) {
	argv := new(api.AuthorizeCheckReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("AuthorizeCheck parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AuthorizeCheck request, IP=%v", httputil.IP(r))
	if argv.ClientId == "" {
		log.Info("%s: missing argument client_id", argv.CommandName())
		svr.responseErrorCode(w, api.ErrorCode_MissingArgument, "missing client_id")
		return
	}
	client, err := svr.clientRepo.GetClient(argv.ClientId)
	if err != nil {
		log.Error("%s: get client %s error: %v", argv.CommandName(), argv.ClientId, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if client == nil {
		log.Info("%s: client %s not found", argv.CommandName(), argv.ClientId)
		svr.responseErrorCode(w, api.ErrorCode_ClientNotFound, "client "+argv.ClientId+" not found")
		return
	}
	session := svr.getSession(r)
	if session == nil {
		log.Info("%s: session not found", argv.CommandName())
		svr.responseErrorCode(w, api.ErrorCode_SessionNotFound, "session-not-found")
		return
	}
	user, err := svr.userRepo.GetUser(session.Uid)
	if err != nil {
		log.Error("%s: get user %d error: %v", argv.CommandName(), session.Uid, err)
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		log.Info("%s: user %d not found", argv.CommandName(), session.Uid)
		svr.responseErrorCode(w, api.ErrorCode_UserNotFound, "user-not-found")
		return
	}
	svr.response(w, http.StatusOK, api.AuthorizeCheckRes{
		Application: client.Name,
		Username:    user.Nickname,
	})
}
