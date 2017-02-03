package model

type AccountType int

// 账号类型
const (
	AccountType_Normal   AccountType = 0 // 大小写字母开头,只允许字母,数字和下划线,长度2~20
	AccountType_Mobile   AccountType = 1 // 手机号
	AccountType_Email    AccountType = 2 // 邮箱
	AccountType_QQ       AccountType = 3 // 从此往下均为第三方账号,所有第三方账号的account前面均加 <type># 前缀
	AccountType_WeChat   AccountType = 4 // 比如 QQ 的加前缀 3#
	AccountType_Github   AccountType = 5 // 微信的加前缀 4#
	AccountType_Facebook AccountType = 6 // 依此类推...
	AccountType_Twitter  AccountType = 7 // 这样做的目的在于确保 account 字段的唯一性

)
