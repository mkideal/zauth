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
		log.Warn("AuthorizeCheck parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AuthorizeCheck request, IP=%v", httputil.IP(r))
	if argv.ClientId == "" {
		svr.responseErrorCode(w, api.ErrorCode_MissingArgument, "missing client_id")
		return
	}
	client, err := svr.clientRepo.FindClient(argv.ClientId)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if client == nil {
		svr.responseErrorCode(w, api.ErrorCode_ClientNotFound, "client "+argv.ClientId+" not found")
		return
	}
	session := svr.getSession(r)
	if session == nil {
		svr.responseErrorCode(w, api.ErrorCode_SessionNotFound, "session-not-found")
		return
	}
	user, err := svr.userRepo.FindUser(session.Uid)
	if err != nil {
		svr.response(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		svr.responseErrorCode(w, api.ErrorCode_UserNotFound, "user-not-found")
		return
	}
	svr.response(w, http.StatusOK, api.AuthorizeCheckRes{
		Application: client.Name,
		Username:    user.Nickname,
	})
}
