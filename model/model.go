package model

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidateClient(client *Client, clientSecret string) bool {
	return clientSecret == client.Secret
}

func ValidatePassword(user *User, password string) bool {
	return user.EncryptedPassword == EncryptPassword(password, user.PasswordSalt)
}

func EncryptPassword(password, salt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+":"+salt)))
}

func JoinAccount(accountType AccountType, account string) string {
	if accountType <= AccountType_Email {
		return account
	}
	return fmt.Sprintf("%d#%s", accountType, account)
}

func SplitAccount(account string) (AccountType, string) {
	var (
		accountType = AccountType_Normal
		index       = strings.Index(account, "#")
	)
	if index > 0 {
		t, err := strconv.Atoi(account[:index])
		if err != nil {
			return accountType, ""
		}
		accountType = AccountType(t)
		account = account[index+1:]
	} else {
		if IsNormalUsername(account) {
			accountType = AccountType_Normal
		} else if IsAutoUsername(account) {
			accountType = AccountType_Auto
		} else if IsEmail(account) {
			accountType = AccountType_Email
		} else if IsTelno(account) {
			accountType = AccountType_Telno
		} else {
			account = ""
		}
	}
	return accountType, account
}

var (
	regNormalUsername = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{1,31}$")
	regAutoUsername   = regexp.MustCompile("^[_][a-zA-Z][a-zA-Z0-9]{1,30}$")
	regEmail          = regexp.MustCompile("^[a-zA-Z0-9_-]{1,64}@.+\\..+$")
	regTelno          = regexp.MustCompile("^[+]?[0-9]{6,15}$")
	regPassword       = regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*_]{6,32}$")
)

func IsNormalUsername(account string) bool { return regNormalUsername.MatchString(account) }
func IsAutoUsername(account string) bool   { return regAutoUsername.MatchString(account) }
func IsEmail(account string) bool          { return regEmail.MatchString(account) }
func IsTelno(account string) bool          { return regTelno.MatchString(account) }
func IsPassword(password string) bool      { return regPassword.MatchString(password) }

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

func ParseTime(s string) (time.Time, error) {
	return time.Parse(RFC3339Milli, s)
}

func ToUnix(s string) int64 {
	t, err := ParseTime(s)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func ToUnixMilli(s string) int64 {
	t, err := ParseTime(s)
	if err != nil {
		return 0
	}
	return t.UnixNano() / 1000000
}

func FormatTime(t time.Time) string {
	return t.Format(RFC3339Milli)
}

func DurationFrom(s string, from time.Time) time.Duration {
	to, err := ParseTime(s)
	if err != nil {
		return 0
	}
	d := to.Sub(from)
	if d < 0 {
		d = 0
	}
	return d
}

func DurationFromNow(s string) time.Duration {
	return DurationFrom(s, time.Now())
}

func IsExpired(s string) bool {
	t, err := ParseTime(s)
	if err != nil {
		return true
	}
	return t.Before(time.Now())
}
