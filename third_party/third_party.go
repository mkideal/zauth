package third_party

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mkideal/accountd/model"
)

// client names
const (
	Facebook = "facebook"
	Github   = "github"
	Twitter  = "twitter"
	WeChat   = "wechat"
	QQ       = "qq"
)

var typesMap = map[string]model.AccountType{
	Facebook: model.AccountType_Facebook,
	Github:   model.AccountType_Github,
	Twitter:  model.AccountType_Twitter,
	WeChat:   model.AccountType_WeChat,
	QQ:       model.AccountType_QQ,
}

var revTypesMap = make(map[model.AccountType]string)

func init() {
	for k, v := range typesMap {
		revTypesMap[v] = k
	}
}

func GetTypeByName(name string) model.AccountType { return typesMap[name] }
func GetNameByType(typ model.AccountType) string  { return revTypesMap[typ] }

var (
	errClientNotFound = errors.New("third party client not found")
)

// error code
const (
	UnsupportedAPI      = "unsupported_api"
	NetworkError        = "network_error"
	ResponseFormatError = "response_format_error"
)

type Error struct {
	Name        string
	Code        string
	Description string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Name, e.Code, e.Description)
}

type AccessTokenResponse struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
	Scope        string
	OpenId       string
	Extra        map[string]string
}

type UserInfoResponse struct {
	OpenId   string
	Nickname string
	Avatar   string
	Sex      model.Gender
	Country  string
	Province string
	City     string
	Extra    map[string]string
}

type Client interface {
	Name() string
	GetAccessToken(clientId, clientSecret, code string) (*AccessTokenResponse, error)
	// openId is optional for many oauth server like github
	GetUserInfo(accessToken, openId string) (*UserInfoResponse, error)
}

type newClient func() Client

var mu sync.Mutex
var clients = make(map[string]newClient)

func Register(name string, f newClient) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, exist := clients[name]; exist {
		panic(fmt.Sprintf("client %s creator existed"))
	}
	clients[name] = f
	return true
}

func New(name string) (Client, error) {
	mu.Lock()
	defer mu.Unlock()
	f, ok := clients[name]
	if !ok {
		return nil, errClientNotFound
	}
	return f(), nil
}
