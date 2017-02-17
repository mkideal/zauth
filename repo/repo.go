package repo

import (
	"bitbucket.org/mkideal/accountd/model"
)

type UserRepository interface {
	AddUser(user *model.User, plainPassword string) error
	UpdateUser(user *model.User) error
	GetUser(uid int64) (*model.User, error)
	GetUserByAccount(account string) (*model.User, error)
	AccountExist(account string) (bool, error)
}

type ClientRepository interface {
	GetClient(clientId string) (*model.Client, error)
}

type AuthorizationRequestRepository interface {
	NewAuthRequest(client *model.Client, uid int64, state, scope, responseType string) (*model.AuthorizationRequest, error)
	GetAuthRequest(clientId, code string) (*model.AuthorizationRequest, error)
	RemoveAuthRequest(id int64) error
}

type TokenRepository interface {
	NewToken(client *model.Client, user *model.User, scope string) (*model.AccessToken, error)
	GetToken(token string) (*model.AccessToken, error)
	RefreshToken(client *model.Client, refreshToken string, scope string) (*model.AccessToken, error)
}

type SessionRepository interface {
	NewSession(uid int64, expireAt string) (*model.Session, error)
	GetSession(sessionId string) (*model.Session, error)
	GetSessionByUid(uid int64) (*model.Session, error)
	UpdateSession(session *model.Session) error
	RemoveSession(sessionId string) error
}

type SqlRepository struct {
}

func (repo SqlRepository) insert(m interface{}) error {
	return nil
}

func (repo SqlRepository) update(m interface{}) error {
	return nil
}

func (repo SqlRepository) remove(m interface{}) error {
	return nil
}

func (repo SqlRepository) get(m interface{}) (bool, error) {
	return false, nil
}

func (repo SqlRepository) getByFields(m interface{}, fields ...string) (bool, error) {
	return false, nil
}

func (repo SqlRepository) has(m interface{}, fields ...string) (bool, error) {
	return false, nil
}
