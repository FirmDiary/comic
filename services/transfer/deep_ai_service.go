package transfer

import (
	"comic/datamodels"
	"comic/repositories"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	oldFix = 1

	api    = "https://api.deepai.org/api/colorizer"
	apiKey = "764f9fc0-97de-4fe3-bf26-0c4ee9052139"
)

type IDeepAiService interface {
	TransferOldFix(file multipart.File, userId int64) (filename string, direction string, err error)
}

type DeepAiService struct {
	repository repositories.IUploadRepository
}

func NewDeepAiService() IDeepAiService {
	return &DeepAiService{repository: repositories.NewUploadRepository()}
}

func (d DeepAiService) TransferOldFix(file multipart.File, userId int64) (filename string, direction string, err error) {
	filename = SaveImgFileToLocal(file, In)
	fileUrlFull := GetFileUrl(filename, In)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api, strings.NewReader("image="+fileUrlFull))
	if err != nil {
		panic(err)
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(resp)
		fmt.Println(string(body))
		return filename, direction, errors.New("解析出现错误，请重试")
	}
	cc, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(resp)
		fmt.Println(string(body))
		return filename, direction, errors.New("json解析出现错误，请重试")
	}

	resUrl := cc.Get("output_url").MustString()
	_, direction = SaveImgUrlToLocal(resUrl, filename, Out)

	//添加数据库记录
	d.repository.Create(&datamodels.Upload{
		File:   filename,
		UserId: userId,
		Type:   oldFix,
	})

	filename += ImgType
	return
}
