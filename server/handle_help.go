package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleHelp(w http.ResponseWriter, r *http.Request) {
	argv, err := parseHelp(r)
	if err != nil {
		log.Warn("Help parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("Help request, IP=%v", httputil.IP(r))
}
