package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	argv, err := parseUser(r)
	if err != nil {
		log.Warn("User parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("User request, IP=%v", httputil.IP(r))
}
