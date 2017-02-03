package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	argv, err := parseSignup(r)
	if err != nil {
		log.Warn("Signup parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("Signup request, IP=%v", httputil.IP(r))
}
