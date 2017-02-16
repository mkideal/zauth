package repo

import (
	"bitbucket.org/mkideal/accountd/model"
)

type UserRepository interface {
	AddUser(user *model.User) error
	UpdateUser(user *model.User) (int, error)
	RemoveUser(uid int64) (int, error)
	FindUser(uid int64) (*model.User, error)
	FindUserByAccount(account string) (*model.User, error)
	AccountExist(account string) (bool, error)
}

type ClientRepository interface {
	FindClient(clientId string) (*model.Client, error)
}

type AuthorizationRequestRepository interface {
	NewAuthRequest(clientId string, uid int64, state, scope string) (*model.AuthorizationRequest, error)
	FindAuthRequest(clientId, code string) (*model.AuthorizationRequest, error)
	RemoveAuthRequest(id int64) error
}

type TokenRepository interface {
	NewToken(client *model.Client, user *model.User, scope string) (*model.AccessToken, error)
	FindToken(token string) (*model.AccessToken, error)
	RefreshToken(client *model.Client, refreshToken string, scope string) (*model.AccessToken, error)
}

type SessionRepository interface {
	NewSession(uid int64, expireAt string) (*model.Session, error)
	FindSession(sessionId string) (*model.Session, error)
	FindSessionByUid(uid int64) (*model.Session, error)
	UpdateSession(session *model.Session) error
	RemoveSession(sessionId string) error
}

type sqlRepository struct {
}
