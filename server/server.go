package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/netutil"
	"github.com/mkideal/pkg/netutil/httputil"

	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/oauth2"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Port                  uint16 `cli:"p,port" usage:"HTTP port" dft:"5200"`
	Mode                  string `cli:"m,mode" usage:"run mode: debug/release" dft:"release"`
	CookieKey             string `cli:"cookie" usage:"cookie key" dft:"accountd"`
	SessionExpireDuration int64  `cli:"session-expire-duration" usage:"session expire duration" dft:"3600"`

	Pages
}

type Pages struct {
	Authorize string `cli:"page-authorize" usage:"web page URL for authorize" dft:"/authorize.html"`
	Login     string `cli:"page-login" usage:"web page URL for login" dft:"/login.html"`
}

type Server struct {
	config Config

	// repositories
	userRepo    UserRepository
	clientRepo  ClientRepository
	authRepo    AuthorizationRequestRepository
	tokenRepo   TokenRepository
	sessionRepo SessionRepository

	running int32
}

func New(config Config) *Server {
	svr := &Server{
		config: config,
	}
	// TODO: initialize repositories
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

// HTTP response methods

func (svr *Server) oauthErrorResponse(cmd string, w http.ResponseWriter, code string, descriptions ...string) error {
	description := ""
	if len(descriptions) > 0 {
		description = strings.Join(descriptions, ". ")
	}
	err := oauth2.NewError(code, description)
	log.Warn("%s oauth error: %v", cmd, err)
	return svr.response(w, http.StatusBadRequest, err)
}

func (svr *Server) errorResponse(cmd string, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if oauthErr, ok := err.(oauth2.Error); ok {
		svr.response(w, http.StatusBadRequest, oauthErr)
	} else {
		svr.oauthErrorResponse(cmd, w, oauth2.ErrorServerError, err.Error())
	}
}

func (svr *Server) response(w http.ResponseWriter, status int, v interface{}) error {
	return httputil.JSONResponse(w, status, v, svr.config.Mode == Debug)
}

func (svr *Server) getTokenFromHeader(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	bearer := oauth2.TokenHeaderPrefix
	if strings.HasPrefix(authorization, bearer) {
		return strings.TrimPrefix(authorization, bearer)
	}
	return ""
}

// get and set session

func (svr *Server) getSession(r *http.Request) *model.Session {
	cookie := r.FormValue("cookie")
	if cookie == "" {
		c, err := r.Cookie(svr.config.CookieKey)
		if err != nil || c == nil {
			return nil
		}
		cookie = c.Value
	}
	log.Trace("cookie: %s", cookie)
	session, _ := svr.sessionRepo.FindSession(cookie)
	return session
}

func (svr *Server) setSession(w http.ResponseWriter, r *http.Request, uid int64) (session *model.Session, err error) {
	session = svr.getSession(r)
	if session != nil {
		svr.sessionRepo.SetSessionUserId(session.Id, uid)
	} else {
		session, err = svr.sessionRepo.NewSession(uid)
		if err != nil {
			return
		}
	}
	cookie := &http.Cookie{
		Name:    svr.config.CookieKey,
		Value:   session.Id,
		Expires: time.Unix(session.Expire, 0),
		MaxAge:  int(svr.config.SessionExpireDuration),
	}
	log.Debug("SetCookie %s for user %d", session.Id, uid)
	http.SetCookie(w, cookie)
	return
}
