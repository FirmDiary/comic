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

//func (m *Mini) ImgSecCheck(url string) (res auth.ResCode2Session, err error) {
//    token, err := m.app.GetAuth().AccessTokenHandle.GetAccessToken()
//
//    var obj interface{}
//
//    bodyBuf := &bytes.Buffer{}
//    bodyWriter := multipart.NewWriter(bodyBuf)
//
//    fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
//    if e != nil {
//        err = fmt.Errorf("error writing to buffer , err=%v", e)
//        return
//    }
//
//    fh, e := os.Open(field.Filename)
//    if e != nil {
//        err = fmt.Errorf("error opening file , err=%v", e)
//        return
//    }
//    defer fh.Close()
//
//    if _, err = io.Copy(fileWriter, fh); err != nil {
//        return
//    }
//
//    contentType := bodyWriter.FormDataContentType()
//    bodyWriter.Close()
//
//    resp, e := http.Post(uri, contentType, bodyBuf)
//    if e != nil {
//        err = e
//        return
//    }
//    defer resp.Body.Close()
//    if resp.StatusCode != http.StatusOK {
//        return nil, err
//    }
//    respBody, err = ioutil.ReadAll(resp.Body)
//    return
//
//    response, err = util.PostJSON("https://api.weixin.qq.com/wxa/img_sec_check?access_token="+token, obj{
//        "access_token": token,
//        "media":
//    })
//
//    return
//}
