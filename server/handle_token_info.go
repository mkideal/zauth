package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleTokenInfo(w http.ResponseWriter, r *http.Request) {
	argv, err := parseTokenInfo(r)
	if err != nil {
		log.Warn("TokenInfo parse arguments error: %v, IP=%v", err, httputil.IP(r))
		return
	}
	log.WithJSON(argv).Debug("TokenInfo request, IP=%v", httputil.IP(r))
}
