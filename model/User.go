// NOTE: AUTO-GENERATED by midc, DON'T edit!!

package model

import (
	"fmt"

	"github.com/mkideal/pkg/storage"
	"github.com/mkideal/pkg/typeconv"
	"gopkg.in/redis.v5"
)

var (
	_ = fmt.Printf
	_ = storage.Unused
	_ = typeconv.Unused
	_ = redis.Nil
)

// Table

// 用户信息
type User struct {
	Id                int64       `xorm:"pk BIGINT(20)"`       // 随机唯一Id
	AccountType       AccountType `xorm:"TEXT"`                // 账号类型
	Account           string      `xorm:"VARCHAR(128) UNIQUE"` // 账号
	Nickname          string      `xorm:"VARCHAR(64)"`         // 昵称
	Avatar            string      `xorm:"VARCHAR(256)"`        // 头像
	Country           string      `xorm:"VARCHAR(32)"`         // 国家
	Province          string      `xorm:"VARCHAR(64)"`         // 省
	City              string      `xorm:"VARCHAR(256)"`        // 城市
	Gender            Gender      `xorm:"TEXT"`                // 性别
	Birthday          string      `xorm:"VARCHAR(32)"`         // 生日
	IdCardType        IdCardType  `xorm:"TEXT"`                // 身份证件类型
	IdCard            string      `xorm:"VARCHAR(64)"`         // 证件唯一标识
	EncryptedPassword string      `xorm:"VARCHAR(64)"`         // 加密后密码
	PasswordSalt      string      `xorm:"VARCHAR(32)"`         // 加密密码的盐
	CreatedAt         string      `xorm:"VARCHAR(32)"`         // 账号创建时间
	CreatedIp         string      `xorm:"VARCHAR(32)"`         // 账号创建时IP
	LastLoginAt       string      `xorm:"VARCHAR(32)"`         // 最后登陆时间
	LastLoginIp       string      `xorm:"VARCHAR(32)"`         // 最后登陆时IP

}

func NewUser() *User {
	return &User{}
}

func (User) Meta() UserMeta               { return userMetaVar }
func (User) TableMeta() storage.TableMeta { return userMetaVar }
func (x User) Key() interface{}           { return x.Id }
func (x *User) SetKey(value string) error {
	return typeconv.String2Int64(&x.Id, value)
}

func (x User) GetField(field string) (interface{}, bool) {
	switch field {
	case userMetaVar.F_account_type:
		return x.AccountType, true
	case userMetaVar.F_account:
		return x.Account, true
	case userMetaVar.F_nickname:
		return x.Nickname, true
	case userMetaVar.F_avatar:
		return x.Avatar, true
	case userMetaVar.F_country:
		return x.Country, true
	case userMetaVar.F_province:
		return x.Province, true
	case userMetaVar.F_city:
		return x.City, true
	case userMetaVar.F_gender:
		return x.Gender, true
	case userMetaVar.F_birthday:
		return x.Birthday, true
	case userMetaVar.F_id_card_type:
		return x.IdCardType, true
	case userMetaVar.F_id_card:
		return x.IdCard, true
	case userMetaVar.F_encrypted_password:
		return x.EncryptedPassword, true
	case userMetaVar.F_password_salt:
		return x.PasswordSalt, true
	case userMetaVar.F_created_at:
		return x.CreatedAt, true
	case userMetaVar.F_created_ip:
		return x.CreatedIp, true
	case userMetaVar.F_last_login_at:
		return x.LastLoginAt, true
	case userMetaVar.F_last_login_ip:
		return x.LastLoginIp, true
	}
	return nil, false
}

func (x *User) SetField(field, value string) error {
	switch field {
	case userMetaVar.F_account_type:
		if err := typeconv.String2Object(&x.AccountType, value); err != nil {
			return err
		}
	case userMetaVar.F_account:
		x.Account = value
	case userMetaVar.F_nickname:
		x.Nickname = value
	case userMetaVar.F_avatar:
		x.Avatar = value
	case userMetaVar.F_country:
		x.Country = value
	case userMetaVar.F_province:
		x.Province = value
	case userMetaVar.F_city:
		x.City = value
	case userMetaVar.F_gender:
		if err := typeconv.String2Object(&x.Gender, value); err != nil {
			return err
		}
	case userMetaVar.F_birthday:
		x.Birthday = value
	case userMetaVar.F_id_card_type:
		if err := typeconv.String2Object(&x.IdCardType, value); err != nil {
			return err
		}
	case userMetaVar.F_id_card:
		x.IdCard = value
	case userMetaVar.F_encrypted_password:
		x.EncryptedPassword = value
	case userMetaVar.F_password_salt:
		x.PasswordSalt = value
	case userMetaVar.F_created_at:
		x.CreatedAt = value
	case userMetaVar.F_created_ip:
		x.CreatedIp = value
	case userMetaVar.F_last_login_at:
		x.LastLoginAt = value
	case userMetaVar.F_last_login_ip:
		x.LastLoginIp = value
	}
	return nil
}

// Meta
type UserMeta struct {
	F_account_type       string
	F_account            string
	F_nickname           string
	F_avatar             string
	F_country            string
	F_province           string
	F_city               string
	F_gender             string
	F_birthday           string
	F_id_card_type       string
	F_id_card            string
	F_encrypted_password string
	F_password_salt      string
	F_created_at         string
	F_created_ip         string
	F_last_login_at      string
	F_last_login_ip      string
}

func (UserMeta) Name() string     { return "user" }
func (UserMeta) Key() string      { return "id" }
func (UserMeta) Fields() []string { return _user_fields }

var userMetaVar = UserMeta{
	F_account_type:       "account_type",
	F_account:            "account",
	F_nickname:           "nickname",
	F_avatar:             "avatar",
	F_country:            "country",
	F_province:           "province",
	F_city:               "city",
	F_gender:             "gender",
	F_birthday:           "birthday",
	F_id_card_type:       "id_card_type",
	F_id_card:            "id_card",
	F_encrypted_password: "encrypted_password",
	F_password_salt:      "password_salt",
	F_created_at:         "created_at",
	F_created_ip:         "created_ip",
	F_last_login_at:      "last_login_at",
	F_last_login_ip:      "last_login_ip",
}

var _user_fields = []string{
	userMetaVar.F_account_type,
	userMetaVar.F_account,
	userMetaVar.F_nickname,
	userMetaVar.F_avatar,
	userMetaVar.F_country,
	userMetaVar.F_province,
	userMetaVar.F_city,
	userMetaVar.F_gender,
	userMetaVar.F_birthday,
	userMetaVar.F_id_card_type,
	userMetaVar.F_id_card,
	userMetaVar.F_encrypted_password,
	userMetaVar.F_password_salt,
	userMetaVar.F_created_at,
	userMetaVar.F_created_ip,
	userMetaVar.F_last_login_at,
	userMetaVar.F_last_login_ip,
}

// Slice
type UserSlice []User

func NewUserSlice(cap int) *UserSlice {
	s := UserSlice(make([]User, 0, cap))
	return &s
}

func (s UserSlice) TableMeta() storage.TableMeta { return userMetaVar }
func (s UserSlice) Len() int                     { return len(s) }
func (s *UserSlice) Slice() []User               { return []User(*s) }

func (s *UserSlice) New(table string, index int, key string) (storage.Table, error) {
	for len(*s) <= index {
		*s = append(*s, User{})
	}
	x := &((*s)[index])
	err := x.SetKey(key)
	return x, err
}

// View
type UserView struct {
	User
}

type UserViewSlice []UserView

func NewUserViewSlice(cap int) *UserViewSlice {
	s := UserViewSlice(make([]UserView, 0, cap))
	return &s
}

func (s UserViewSlice) TableMeta() storage.TableMeta { return userMetaVar }
func (s UserViewSlice) Len() int                     { return len(s) }
func (s *UserViewSlice) Slice() []UserView           { return []UserView(*s) }

func (s *UserViewSlice) New(table string, index int, key string) (storage.Table, error) {
	if table == "user" {
		for len(*s) <= index {
			x := User{}
			*s = append(*s, UserView{User: x})
		}
		x := &((*s)[index].User)
		err := x.SetKey(key)
		return x, err
	}
	v := &((*s)[index])
	for t, x := range v.tables() {
		if t == table {
			err := x.SetKey(key)
			return x, err
		}
	}
	return nil, storage.ErrTableNotFoundInView
}

var (
	UserViewVar  = UserView{}
	userViewRefs = map[string]storage.View{}
)

func (UserView) TableMeta() storage.TableMeta  { return userMetaVar }
func (UserView) Fields() storage.FieldList     { return storage.FieldSlice(userMetaVar.Fields()) }
func (UserView) Refs() map[string]storage.View { return userViewRefs }
func (view *UserView) tables() map[string]storage.Table {
	m := make(map[string]storage.Table)
	m["user"] = &view.User
	return m
}
