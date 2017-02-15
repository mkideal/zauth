package model

type AccountType int

// 账号类型
const (
	AccountType_Normal   AccountType = 0 // 大小写字母开头,只允许字母,数字和下划线,长度2~20
	AccountType_Auto     AccountType = 1 // 自动生成账号(快捷注册),account 就是uid
	AccountType_Telno    AccountType = 2 // 手机号
	AccountType_Email    AccountType = 3 // 邮箱
	AccountType_QQ       AccountType = 4 // 从此往下均为第三方账号,所有第三方账号的account前面均加 <type># 前缀
	AccountType_WeChat   AccountType = 5 // 比如 QQ 的加前缀 4#
	AccountType_Github   AccountType = 6 // 微信的加前缀 5#
	AccountType_Facebook AccountType = 7 // 依此类推...
	AccountType_Twitter  AccountType = 8 // 这样做的目的在于确保 account 字段的唯一性

)
