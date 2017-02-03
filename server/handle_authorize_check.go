package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
)

func (svr *Server) handleAuthorizeCheck(w http.ResponseWriter, r *http.Request) {
	argv, err := parseAuthorizeCheck(r)
	if err != nil {
		log.Warn("AuthorizeCheck parse arguments error: %v, IP=%v", err, httputil.IP(r))
		svr.response(w, http.StatusBadRequest, err)
		return
	}
	log.WithJSON(argv).Debug("AuthorizeCheck request, IP=%v", httputil.IP(r))
}
