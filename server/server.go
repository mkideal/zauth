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

	_ "github.com/go-sql-driver/mysql"

	"bitbucket.org/mkideal/accountd/api"
	"bitbucket.org/mkideal/accountd/model"
	"bitbucket.org/mkideal/accountd/oauth2"
	"bitbucket.org/mkideal/accountd/repo"
)

const (
	Debug   = "debug"
	Release = "release"
)

type Config struct {
	Driver                string `cli:"driver" usage:"sql database driver: mysql" dft:"mysql"`
	DataSourceName        string `cli:"dsn" usage:"data source name for specified driver" dft:"accountd:aXXpwd@/accountdb"`
	Port                  uint16 `cli:"p,port" usage:"HTTP port" dft:"5200"`
	Mode                  string `cli:"m,mode" usage:"running mode: debug/release" dft:"release"`
	CookieKey             string `cli:"cookie" usage:"cookie key" dft:"accountd"`
	SessionExpireDuration int64  `cli:"session-expire-duration" usage:"session expire duration(seconds)" dft:"3600"`
	HTMLDir               string `cli:"html" usage:"HTML static directory" dft:"html"`
	HTMLRoouter           string `cli:"html-router" usage:"HTML static files router" dft:"/"`

	Pages
}

type Pages struct {
	Authorize string `cli:"page-authorize" usage:"web page URL for authorize" dft:"/authorize.html"`
	Login     string `cli:"page-login" usage:"web page URL for login" dft:"/login.html"`
}

type Server struct {
	config Config

	// repositories
	userRepo    repo.UserRepository
	clientRepo  repo.ClientRepository
	authRepo    repo.AuthorizationRequestRepository
	tokenRepo   repo.TokenRepository
	sessionRepo repo.SessionRepository

	running int32
}

func New(config Config) (*Server, error) {
	sqlRepo, err := repo.NewSqlRepository(config.Driver, config.DataSourceName)
	if err != nil {
		return nil, err
	}

	svr := &Server{
		config: config,
	}
	// initialize repositories
	svr.userRepo = repo.NewUserRepository(sqlRepo)
	svr.clientRepo = repo.NewClientRepository(sqlRepo)
	svr.authRepo = repo.NewAuthorizationRequestRepository(sqlRepo)
	svr.tokenRepo = repo.NewTokenRepository(sqlRepo)
	svr.sessionRepo = repo.NewSessionRepository(sqlRepo)
	return svr, nil
}

func (svr *Server) registerHandler(mux *httputil.ServeMux, pattern, method string, h http.HandlerFunc) {
	mux.Handle(pattern, httputil.NewHandlerFunc(method, h))
}

func (svr *Server) Run() error {
	if !atomic.CompareAndSwapInt32(&svr.running, 0, 1) {
		return fmt.Errorf("server rerun")
	}

	log.WithJSON(svr.config).Info("config of server")

	// register HTTP api
	mux := httputil.NewServeMux()
	svr.registerAllHandlers(mux)
	mux.Handle(svr.config.HTMLRoouter, http.FileServer(http.Dir(svr.config.HTMLDir)))

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
}

// HTTP response methods

func (svr *Server) response(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return httputil.Response(w, http.StatusOK, r.Header.Get(httputil.HeaderAccept), v, svr.config.Mode == Debug)
}

func (svr *Server) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	debug := svr.config.Mode == Debug
	switch e := err.(type) {
	case oauth2.OAuthErrorCode:
		e2 := api.NewError(e.Error(), "")
		httputil.Response(w, e2.Status(), r.Header.Get(httputil.HeaderAccept), e2, debug)
	case api.Error:
		httputil.Response(w, e.Status(), r.Header.Get(httputil.HeaderAccept), e, debug)
	case api.ErrorCode:
		httputil.Response(w, e.Status(), r.Header.Get(httputil.HeaderAccept), e.NewError(""), debug)
	default:
		e2 := api.ErrorCode_InternalServerError.NewError(e.Error())
		httputil.Response(w, e2.Status(), r.Header.Get(httputil.HeaderAccept), e2, debug)
	}
}

// get and set token/session

func (svr *Server) getTokenFromHeader(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if strings.HasPrefix(authorization, oauth2.TokenHeaderPrefix) {
		return strings.TrimPrefix(authorization, oauth2.TokenHeaderPrefix)
	}
	return ""
}

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
	session, _ := svr.sessionRepo.GetSession(cookie)
	return session
}

func (svr *Server) setSession(w http.ResponseWriter, r *http.Request, uid int64) (session *model.Session, err error) {
	session = svr.getSession(r)
	duration := time.Hour * 24 * 3650 // 10 years
	if svr.config.SessionExpireDuration > 0 {
		duration = time.Duration(svr.config.SessionExpireDuration) * time.Second
	}
	expireAt := time.Now().Add(duration)
	if session != nil {
		session.Uid = uid
		session.ExpireAt = model.FormatTime(expireAt)
		if err = svr.sessionRepo.UpdateSession(session); err != nil {
			return
		}
	} else {
		session, err = svr.sessionRepo.NewSession(uid, model.FormatTime(expireAt))
		if err != nil {
			return
		}
	}
	cookie := &http.Cookie{
		Name:    svr.config.CookieKey,
		Value:   session.Id,
		Expires: expireAt,
		MaxAge:  int(duration / time.Second),
	}
	log.Debug("SetCookie %s for user %d", session.Id, uid)
	http.SetCookie(w, cookie)
	return
}

// authorization client
func (svr *Server) clientAuth(cmd string, w http.ResponseWriter, r *http.Request) *model.Client {
	clientId, clientSecret, ok := r.BasicAuth()
	if !ok {
		log.Info("%s: Client BasicAuth failed", cmd)
		svr.errorResponse(w, r, api.ErrorCode_ClientUnauthorized.NewError("client BasicAuth failed"))
		return nil
	}
	client, err := svr.clientRepo.GetClient(clientId)
	if err != nil {
		log.Error("%s: GetClient %s error: %v", cmd, clientId, err)
		svr.errorResponse(w, r, err)
		return nil
	}
	if client == nil {
		log.Info("%s: Client %s not found", cmd, clientId)
		svr.errorResponse(w, r, api.ErrorCode_ClientNotFound)
		return nil
	}
	if !model.ValidateClint(client, clientSecret) {
		log.Info("%s: Client %s secret invalid", cmd, clientId)
		svr.errorResponse(w, r, api.ErrorCode_IncorrectClientSecret)
		return nil
	}
	return client
}

func makeUserInfo(user *model.User) api.UserInfo {
	return api.UserInfo{
		Id:          user.Id,
		Account:     user.Account,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Qrcode:      user.Qrcode,
		Gender:      int(user.Gender),
		Birthday:    user.Birthday,
		LastLoginAt: user.LastLoginAt,
		LastLoginIp: user.LastLoginIp,
	}
}

func makeTokenInfo(token *model.Token) api.TokenInfo {
	return api.TokenInfo{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Scope:        token.Scope,
		ExpireAt:     token.ExpireAt,
	}
}

func makeRouter(command string) string {
	return "/v1/" + command
}
