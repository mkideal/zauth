package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	argv, err := parseLogout(r)
	if err != nil {
		log.Warn("Logout parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("Logout request, IP=%v", httputil.IP(r))
}
