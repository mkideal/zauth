package server

import (
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil/httputil"
	"github.com/mkideal/pkg/textutil/namemapper"

	"github.com/mkideal/accountd/api"
)

func (svr *Server) handleHelp(w http.ResponseWriter, r *http.Request) {
	ip := httputil.IP(r)
	argv := new(api.HelpReq)
	err := argv.Parse(r)
	if err != nil {
		log.Info("Help parse arguments error: %v, IP=%v", err, ip)
		svr.errorResponse(w, r, api.ErrorCode_BadArgument.NewError(err.Error()))
		return
	}
	log.WithJSON(argv).Debug("Help request, IP=%v", ip)
	if argv.Cmd == "" {
		type response struct {
			Commands []string `json:"commands"`
			Routers  []string `json:"routers"`
		}
		res := new(response)
		commands := api.Commands()
		for _, cmd := range commands {
			res.Commands = append(res.Commands, cmd.CommandName())
			res.Routers = append(res.Routers, cmd.CommandMethod()+" "+makeRouter(namemapper.UnderScore(cmd.CommandName())))
		}
		svr.response(w, r, res)
	} else {
		req := api.GetCommand(argv.Cmd)
		if req == nil {
			svr.errorResponse(w, r, api.ErrorCode_CommandNotFound.NewError("command "+argv.Cmd+" not found"))
			return
		}
		type response struct {
			Command   string      `json:"command"`
			Method    string      `json:"method"`
			Router    string      `json:"router"`
			Arguments interface{} `json:"arguments"`
		}
		svr.response(w, r, response{
			Command:   req.CommandName(),
			Method:    req.CommandMethod(),
			Router:    makeRouter(namemapper.UnderScore(req.CommandName())),
			Arguments: req,
		})
	}
}
