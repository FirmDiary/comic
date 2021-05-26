package wechat

import (
    "errors"
    "github.com/silenceper/wechat/v2/miniprogram/auth"
    "github.com/sirupsen/logrus"
)

func (m *Mini) GetAccessToken() (ak string, err error) {
    ak, err = m.app.GetAuth().GetAccessToken()
    if err != nil {
        logrus.Error(err)
        err = errors.New("获取AccessToken失败")
    }
    return
}

func (m *Mini) Code2Session(code string) (res auth.ResCode2Session, err error) {
    return m.app.GetAuth().Code2Session(code)
}
