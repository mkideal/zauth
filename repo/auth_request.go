package repo

import (
	"time"

	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/model"
)

type authorizationRequestRepository struct {
	*SqlRepository
}

func NewAuthorizationRequestRepository(sqlRepo *SqlRepository) AuthorizationRequestRepository {
	return authorizationRequestRepository{SqlRepository: sqlRepo}
}

func (repo authorizationRequestRepository) NewAuthRequest(client *model.Client, uid int64, state, scope, responseType string) (*model.AuthorizationRequest, error) {
	ar := &model.AuthorizationRequest{
		CreatedAt:         model.FormatTime(time.Now()),
		AuthorizationCode: random.String(64, random.CryptoSource, random.O_DIGIT, random.O_UPPER_CHAR, random.O_LOWER_CHAR),
		ClientId:          client.Id,
		Uid:               uid,
		State:             state,
		RedirectURI:       client.CallbackURL,
		ResponseType:      responseType,
	}
	err := repo.Insert(ar)
	if err != nil {
		ar = nil
	}
	return ar, err
}

func (repo authorizationRequestRepository) GetAuthRequest(clientId, code string) (*model.AuthorizationRequest, error) {
	ar := &model.AuthorizationRequest{ClientId: clientId, AuthorizationCode: code}
	found, err := repo.Get(ar)
	if !found || err != nil {
		ar = nil
	}
	return ar, err
}

func (repo authorizationRequestRepository) RemoveAuthRequest(code string) error {
	return repo.Remove(&model.AuthorizationRequest{AuthorizationCode: code})
}
