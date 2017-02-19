package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleSignout(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.SignoutReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Signout parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("Signout request, IP=%v", ip)
	session := svr.getSession(r)
	if session != nil {
		if err := svr.sessionRepo.RemoveSession(session.Id); err != nil {
			log.Error("%s: remove session %s error: %v", argv.CommandName(), session.Id, err)
			svr.errorResponse(w, err)
			return
		}
	}
	http.Redirect(w, r, svr.config.Pages.Login, http.StatusFound)
}
