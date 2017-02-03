package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleSignin(w http.ResponseWriter, r *http.Request) {
	argv, err := parseSignin(r)
	if err != nil {
		log.Warn("Signin parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("Signin request, IP=%v", httputil.IP(r))
}
