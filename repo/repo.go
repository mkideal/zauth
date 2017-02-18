package repo

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/midlang/mid/x/go/storage"

	"bitbucket.org/mkideal/accountd/model"
)

type UserRepository interface {
	AddUser(user *model.User, plainPassword string) error
	UpdateUser(user *model.User, fields ...string) error
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
	RemoveAuthRequest(code string) error
}

type TokenRepository interface {
	NewToken(user *model.User, clientId, scope string) (*model.Token, error)
	GetToken(token string) (*model.Token, error)
	RefreshToken(refreshToken string, scope string) (*model.Token, error)
}

type SessionRepository interface {
	NewSession(uid int64, expireAt string) (*model.Session, error)
	GetSession(sessionId string) (*model.Session, error)
	UpdateSession(session *model.Session) error
	RemoveSession(sessionId string) error
}

type SqlRepository struct {
	eng *xorm.Engine
}

func NewSqlRepository(driver, dataSourceName string) (*SqlRepository, error) {
	eng, err := xorm.NewEngine(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &SqlRepository{
		eng: eng,
	}, nil
}

var errWriteDatabaseFailed = errors.New("write to database failed")

func writeOpError(n int64, err error) error {
	if err != nil {
		return err
	}
	if n != 1 {
		return errWriteDatabaseFailed
	}
	return nil
}

func where(m storage.ReadonlyTable, byFields ...string) (query string, values []interface{}, err error) {
	if len(byFields) == 0 {
		query = fmt.Sprintf("%s = ?", m.Meta().Key())
		values = []interface{}{m.Key()}
	} else {
		var buf bytes.Buffer
		for i, field := range byFields {
			if i > 0 {
				buf.WriteString(" and ")
			}
			fmt.Fprintf(&buf, "%s = ?", field)
			value, found := m.GetField(field)
			if !found {
				err = fmt.Errorf("table %s does not contain field %s", m.Meta().Name(), field)
				return
			}
			values = append(values, value)
		}
		query = buf.String()
	}
	return
}

func (repo SqlRepository) Insert(m interface{}) error {
	return writeOpError(repo.eng.InsertOne(m))
}

func (repo SqlRepository) Update(m storage.ReadonlyTable, fields ...string) error {
	if len(fields) == 0 {
		fields = m.Meta().Fields()
	}
	query, values, err := where(m)
	if err != nil {
		return err
	}
	return writeOpError(repo.eng.Cols(fields...).Where(query, values...).Update(m))
}

func (repo SqlRepository) Remove(m storage.ReadonlyTable, byFields ...string) error {
	query, values, err := where(m, byFields...)
	if err != nil {
		return err
	}
	return writeOpError(repo.eng.Where(query, values...).Delete(m))
}

func (repo SqlRepository) Get(m storage.Table, byFields ...string) (bool, error) {
	query, values, err := where(m, byFields...)
	if err != nil {
		return false, err
	}
	return repo.eng.Where(query, values...).Get(m)
}

func (repo SqlRepository) Exist(m storage.ReadonlyTable, byFields ...string) (bool, error) {
	query, values, err := where(m, byFields...)
	if err != nil {
		return false, err
	}
	n, err := repo.eng.Where(query, values...).Count(m)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
