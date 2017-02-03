package server

import (
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil"
	"github.com/mkideal/pkg/netutil/httputil"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Port uint16 `cli:"p,port" usage:"HTTP port" dft:"5200"`
	Mode string `cli:"m,mode" usage:"run mode: debug/release" dft:"release"`
}

type Server struct {
	config Config

	running int32
}

func New(config Config) *Server {
	svr := &Server{
		config: config,
	}
	return svr
}

func (svr *Server) Run() error {
	if !atomic.CompareAndSwapInt32(&svr.running, 0, 1) {
		return fmt.Errorf("server rerun")
	}

	log.WithJSON(svr.config).Info("config of server")

	// register HTTP api
	mux := httputil.NewServeMux()
	svr.registerAllHandlers(mux)

	// listen and serve HTTP service
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", svr.config.Port),
		Handler: mux,
	}
	ln, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		return err
	}
	go httpServer.Serve(netutil.NewTCPKeepAliveListener(ln.(*net.TCPListener), time.Minute*3))

	return nil
}

func (svr *Server) Quit() {
	if !atomic.CompareAndSwapInt32(&svr.running, 1, 0) {
		return
	}
	// TODO
}
