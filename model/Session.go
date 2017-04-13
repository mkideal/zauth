// NOTE: AUTO-GENERATED by midc, DON'T edit!!

package model

import (
	"fmt"

	"github.com/midlang/mid/x/go/storage"
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

// Web 会话
type Session struct {
	Id        string `xorm:"pk VARCHAR(32)"` // 唯一Id,用作cookie
	Uid       int64  `xorm:"BIGINT(20)"`     // 关联的用户Id
	CreatedAt string `xorm:"VARCHAR(32)"`    // 创建时间
	ExpireAt  string `xorm:"VARCHAR(32)"`    // 到期时间

}

func NewSession() *Session {
	return &Session{}
}

func (Session) Meta() SessionMeta            { return sessionMetaVar }
func (Session) TableMeta() storage.TableMeta { return sessionMetaVar }
func (x Session) Key() interface{}           { return x.Id }
func (x *Session) SetKey(value string) error {
	x.Id = value
	return nil
}

func (x Session) GetField(field string) (interface{}, bool) {
	switch field {
	case sessionMetaVar.F_uid:
		return x.Uid, true
	case sessionMetaVar.F_created_at:
		return x.CreatedAt, true
	case sessionMetaVar.F_expire_at:
		return x.ExpireAt, true
	}
	return nil, false
}

func (x *Session) SetField(field, value string) error {
	switch field {
	case sessionMetaVar.F_uid:
		return typeconv.String2Int64(&x.Uid, value)
	case sessionMetaVar.F_created_at:
		x.CreatedAt = value
	case sessionMetaVar.F_expire_at:
		x.ExpireAt = value
	}
	return nil
}

// Meta
type SessionMeta struct {
	F_uid        string
	F_created_at string
	F_expire_at  string
}

func (SessionMeta) Name() string     { return "session" }
func (SessionMeta) Key() string      { return "id" }
func (SessionMeta) Fields() []string { return _session_fields }

var sessionMetaVar = SessionMeta{
	F_uid:        "uid",
	F_created_at: "created_at",
	F_expire_at:  "expire_at",
}

var _session_fields = []string{
	sessionMetaVar.F_uid,
	sessionMetaVar.F_created_at,
	sessionMetaVar.F_expire_at,
}

// Slice
type SessionSlice []Session

func NewSessionSlice(cap int) *SessionSlice {
	s := SessionSlice(make([]Session, 0, cap))
	return &s
}

func (s SessionSlice) TableMeta() storage.TableMeta { return sessionMetaVar }
func (s SessionSlice) Len() int                     { return len(s) }
func (s *SessionSlice) Slice() []Session            { return []Session(*s) }

func (s *SessionSlice) New(table string, index int, key string) (storage.Table, error) {
	for len(*s) <= index {
		*s = append(*s, Session{})
	}
	x := &((*s)[index])
	err := x.SetKey(key)
	return x, err
}

// View
type SessionView struct {
	Session
}

type SessionViewSlice []SessionView

func NewSessionViewSlice(cap int) *SessionViewSlice {
	s := SessionViewSlice(make([]SessionView, 0, cap))
	return &s
}

func (s SessionViewSlice) TableMeta() storage.TableMeta { return sessionMetaVar }
func (s SessionViewSlice) Len() int                     { return len(s) }
func (s *SessionViewSlice) Slice() []SessionView        { return []SessionView(*s) }

func (s *SessionViewSlice) New(table string, index int, key string) (storage.Table, error) {
	if table == "session" {
		for len(*s) <= index {
			x := Session{}
			*s = append(*s, SessionView{Session: x})
		}
		x := &((*s)[index].Session)
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
	SessionViewVar  = SessionView{}
	sessionViewRefs = map[string]storage.View{}
)

func (SessionView) TableMeta() storage.TableMeta  { return sessionMetaVar }
func (SessionView) Fields() storage.FieldList     { return storage.FieldSlice(sessionMetaVar.Fields()) }
func (SessionView) Refs() map[string]storage.View { return sessionViewRefs }
func (view *SessionView) tables() map[string]storage.Table {
	m := make(map[string]storage.Table)
	m["session"] = &view.Session
	return m
}
