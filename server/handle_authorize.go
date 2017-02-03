package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	argv, err := parseAuthorize(r)
	if err != nil {
		log.Warn("Authorize parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("Authorize request, IP=%v", httputil.IP(r))
}
