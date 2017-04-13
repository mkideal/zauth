package repo

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mkideal/pkg/math/random"

	"github.com/mkideal/accountd/model"
)

type telnoVerifyCodeRepository struct {
	*SqlRepository
}

func NewTelnoVerifyCodeRepository(sqlRepo *SqlRepository) TelnoVerifyCodeRepository {
	return &telnoVerifyCodeRepository{SqlRepository: sqlRepo}
}

func (repo *telnoVerifyCodeRepository) NewTelnoCode(codeLength int, telno string, maxIntervalSeconds, expirationSeconds int64) (*model.TelnoVerifyCode, error) {
	vcode := &model.TelnoVerifyCode{Telno: telno}
	found, err := repo.Get(vcode)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if found {
		from, err := model.ParseTime(vcode.CreatedAt)
		if err != nil {
			return nil, err
		}
		if now.Sub(from) < time.Duration(maxIntervalSeconds)*time.Second {
			return nil, nil
		}
	}
	vcode.CreatedAt = model.FormatTime(now)
	vcode.ExpireAt = model.FormatTime(now.Add(time.Second * time.Duration(expirationSeconds)))
	vcode.Code = random.String(codeLength, nil, random.O_DIGIT)
	if found {
		if err := repo.Update(vcode); err != nil {
			return nil, err
		}
	} else {
		if err := repo.Insert(vcode); err != nil {
			return nil, err
		}
	}
	return vcode, nil
}

func (repo *telnoVerifyCodeRepository) FindTelnoCode(telno string) (*model.TelnoVerifyCode, error) {
	vcode := &model.TelnoVerifyCode{Telno: telno}
	found, err := repo.Get(vcode)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return vcode, nil
}

func (repo *telnoVerifyCodeRepository) UpdateTelnoCode(vcode *model.TelnoVerifyCode, fields ...string) error {
	return repo.Update(vcode)
}

func (repo *telnoVerifyCodeRepository) RemoveTelnoCode(telno string) error {
	return repo.Remove(model.TelnoVerifyCode{Telno: telno})
}

func (repo *telnoVerifyCodeRepository) SendTelnoCode(vcode *model.TelnoVerifyCode, uri, un, pw, msgFormat string) error {
	resp, err := http.PostForm(uri, url.Values{
		"un":    {un},
		"pw":    {pw},
		"phone": {vcode.Telno},
		"msg":   {fmt.Sprintf(msgFormat, vcode.Code)},
		"rd":    {"0"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	r := bufio.NewReader(resp.Body)
	line, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	var (
		t      string
		status int
	)
	fmt.Sscanf(string(line), "%s,%d", &t, &status)
	if status != 0 {
		return fmt.Errorf("failed to send SMS: status=%d", status)
	}
	return nil
}
