package repo

import (
	"fmt"
	"time"

	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/model"
)

type userRepository struct {
	SqlRepository
}

func NewUserRepository(sqlRepo SqlRepository) UserRepository {
	return userRepository{SqlRepository: sqlRepo}
}

func (repo userRepository) AddUser(user *model.User, plainPassword string) error {
	// initialize user
	if user.CreatedAt == "" {
		user.CreatedAt = model.FormatTime(time.Now())
	}
	if user.PasswordSalt == "" {
		user.PasswordSalt = random.String(64, nil)
	}
	if user.EncryptedPassword == "" {
		user.EncryptedPassword = model.EncryptPassword(plainPassword, user.PasswordSalt)
	}

	emptyAccount := user.Account == ""
	emptyNickname := user.Nickname == ""
	const maxTryTimes = 3
	for i := 0; i < maxTryTimes; i++ {
		user.Id = random.Int63(random.CryptoSource) % 0xFFFFFFFFFF
		if emptyAccount {
			user.Account = fmt.Sprintf("_%d", user.Id)
		}
		if emptyNickname {
			user.Nickname = user.Account
		}
		if err := repo.insert(user); err == nil {
			break
		} else if i+1 == maxTryTimes {
			return err
		}
	}
	return nil
}

func (repo userRepository) UpdateUser(user *model.User) error {
	return repo.update(user)
}

func (repo userRepository) GetUser(uid int64) (*model.User, error) {
	user := &model.User{Id: uid}
	found, err := repo.get(user)
	if !found || err != nil {
		user = nil
	}
	return user, err
}

func (repo userRepository) GetUserByAccount(account string) (*model.User, error) {
	user := &model.User{Account: account}
	found, err := repo.getByFields(user, model.UserMetaVar.F_account)
	if !found || err != nil {
		user = nil
	}
	return user, err
}

func (repo userRepository) AccountExist(account string) (bool, error) {
	return repo.has(&model.User{Account: account}, model.UserMetaVar.F_account)
}
