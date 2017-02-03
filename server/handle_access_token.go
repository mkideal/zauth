package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleAccessToken(w http.ResponseWriter, r *http.Request) {
	argv, err := parseAccessToken(r)
	if err != nil {
		log.Warn("AccessToken parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("AccessToken request, IP=%v", httputil.IP(r))
}
