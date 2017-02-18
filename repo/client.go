package repo

import (
	"bitbucket.org/mkideal/accountd/model"
)

type clientRepository struct {
	*SqlRepository
}

func NewClientRepository(sqlRepo *SqlRepository) ClientRepository {
	return clientRepository{SqlRepository: sqlRepo}
}

func (repo clientRepository) GetClient(clientId string) (*model.Client, error) {
	client := &model.Client{Id: clientId}
	found, err := repo.Get(client)
	if !found || err != nil {
		client = nil
	}
	return client, err
}
