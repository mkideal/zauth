package repo

import (
	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/model"
)

type sessionRepository struct {
	SqlRepository
}

func NewSessionRepository(sqlRepo SqlRepository) SessionRepository {
	return sessionRepository{SqlRepository: sqlRepo}
}

func (repo sessionRepository) NewSession(uid int64, expireAt string) (*model.Session, error) {
	session := &model.Session{
		Id:       random.String(64, nil),
		Uid:      uid,
		ExpireAt: expireAt,
	}
	err := repo.insert(session)
	if err != nil {
		session = nil
	}
	return session, err
}

func (repo sessionRepository) GetSession(sessionId string) (*model.Session, error) {
	session := &model.Session{Id: sessionId}
	found, err := repo.get(session)
	if !found || err != nil {
		session = nil
	}
	return session, err
}

func (repo sessionRepository) GetSessionByUid(uid int64) (*model.Session, error) {
	session := &model.Session{Uid: uid}
	found, err := repo.getByFields(session, model.SessionMetaVar.F_uid)
	if !found || err != nil {
		session = nil
	}
	return session, err
}

func (repo sessionRepository) UpdateSession(session *model.Session) error {
	return repo.update(session)
}

func (repo sessionRepository) RemoveSession(sessionId string) error {
	return repo.remove(&model.Session{Id: sessionId})
}
