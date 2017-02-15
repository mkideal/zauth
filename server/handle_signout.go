package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleSignout(w http.ResponseWriter, r *http.Request) {
	argv := new(api.SignoutReq)
	err := argv.Parse(r)
	if err != nil {
		log.Warn("Signout parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signout request, IP=%v", httputil.IP(r))
	session := svr.getSession(r)
	if session != nil {
		if err := svr.sessionRepo.RemoveSession(session.Id); err != nil {
			svr.response(w, http.StatusInternalServerError, err)
			return
		}
	}
	http.Redirect(w, r, svr.config.Pages.Login, http.StatusFound)
}
