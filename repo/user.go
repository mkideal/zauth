package repo

import (
	"fmt"
	"time"

	"github.com/mkideal/pkg/math/random"

	"bitbucket.org/mkideal/accountd/model"
)

type userRepository struct {
	*SqlRepository
}

func NewUserRepository(sqlRepo *SqlRepository) UserRepository {
	return userRepository{SqlRepository: sqlRepo}
}

func WithNickname(nickname string) UserAddOption {
	return func(user *model.User) {
		user.Nickname = nickname
	}
}

func WithAccount(account string) UserAddOption {
	return func(user *model.User) {
		user.Account = account
	}
}

func WithGender(gender model.Gender) UserAddOption {
	return func(user *model.User) {
		user.Gender = gender
	}
}

func WithCountry(country string) UserAddOption {
	return func(user *model.User) {
		user.Country = country
	}
}

func WithProvince(province string) UserAddOption {
	return func(user *model.User) {
		user.Province = province
	}
}

func WithCity(city string) UserAddOption {
	return func(user *model.User) {
		user.City = city
	}
}

func (repo userRepository) AddUser(user *model.User, plainPassword string, opts ...UserAddOption) error {
	if plainPassword == "" {
		mode := random.O_DIGIT | random.O_LOWER_CHAR | random.O_UPPER_CHAR
		plainPassword = random.String(16, random.CryptoSource, mode)
	}
	// initialize user
	if user.CreatedAt == "" {
		user.CreatedAt = model.FormatTime(time.Now())
	}
	if user.PasswordSalt == "" {
		user.PasswordSalt = random.String(32, nil)
	}
	if user.EncryptedPassword == "" {
		user.EncryptedPassword = model.EncryptPassword(plainPassword, user.PasswordSalt)
	}
	for _, opt := range opts {
		opt(user)
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
		if err := repo.Insert(user); err == nil {
			break
		} else if i+1 == maxTryTimes {
			return err
		}
	}
	return nil
}

func (repo userRepository) UpdateUser(user *model.User, fields ...string) error {
	return repo.Update(user, fields...)
}

func (repo userRepository) GetUser(uid int64) (*model.User, error) {
	user := &model.User{Id: uid}
	found, err := repo.Get(user)
	if !found || err != nil {
		user = nil
	}
	return user, err
}

func (repo userRepository) GetUserByAccount(account string) (*model.User, error) {
	user := &model.User{Account: account}
	found, err := repo.Get(user, model.UserMetaVar.F_account)
	if !found || err != nil {
		user = nil
	}
	return user, err
}

func (repo userRepository) AccountExist(account string) (bool, error) {
	return repo.Exist(&model.User{Account: account}, model.UserMetaVar.F_account)
}
