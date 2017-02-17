package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/api"
)

func (svr *Server) handleHelp(w http.ResponseWriter, r *http.Request) {
	argv := new(api.HelpReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Help parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Help request, IP=%v", httputil.IP(r))
	svr.response(w, http.StatusOK, "TODO")
}
