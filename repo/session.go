package repo

import (
	"github.com/mkideal/pkg/math/random"

	"github.com/mkideal/accountd/model"
)

type sessionRepository struct {
	*SqlRepository
}

func NewSessionRepository(sqlRepo *SqlRepository) SessionRepository {
	return sessionRepository{SqlRepository: sqlRepo}
}

func (repo sessionRepository) NewSession(uid int64, expireAt string) (*model.Session, error) {
	session := &model.Session{
		Id:       random.String(32, nil),
		Uid:      uid,
		ExpireAt: expireAt,
	}
	err := repo.Insert(session)
	if err != nil {
		session = nil
	}
	return session, err
}

func (repo sessionRepository) GetSession(sessionId string) (*model.Session, error) {
	session := &model.Session{Id: sessionId}
	found, err := repo.Get(session)
	if !found || err != nil {
		session = nil
	}
	return session, err
}

func (repo sessionRepository) UpdateSession(session *model.Session) error {
	return repo.Update(session)
}

func (repo sessionRepository) RemoveSession(sessionId string) error {
	return repo.Remove(&model.Session{Id: sessionId})
}
