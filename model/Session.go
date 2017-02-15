package model

type Session struct {
	Id       string
	Uid      int64
	ExpireAt string
}

type SessionMeta struct {
	F_id       string
	F_uid      string
	F_expireAt string
}

func (SessionMeta) Name() string {
	return "Session"
}

func (SessionMeta) NumField() int {
	return 3
}

func (SessionMeta) Field(i int, v Session) (string, interface{}) {
	switch i {

	case 0:
		return "id", v.Id
	case 1:
		return "uid", v.Uid
	case 2:
		return "expireAt", v.ExpireAt

	}
	return "", nil
}

func (SessionMeta) FieldPtr(i int, v *Session) (string, interface{}) {
	switch i {

	case 0:
		return "id", &v.Id
	case 1:
		return "uid", &v.Uid
	case 2:
		return "expireAt", &v.ExpireAt

	}
	return "", nil
}

var SessionMetaVar = SessionMeta{

	F_id:       "id",
	F_uid:      "uid",
	F_expireAt: "expireAt",
}
