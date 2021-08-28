package services

import (
	"comic/common"
	"comic/datamodels"
	"comic/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	oldFixMT = 3 //美图API老照片修复

	apiOldFixMT = "https://openapi.mtlab.meitu.com/v1/aiquality"

	oldFixMTAppID     = 118725
	oldFixMTAppKey    = "45b1a079c353419d8b2707928b516cf9"
	oldFixMTAppSecret = "d6c320093439457086cd42a8872cde02"
)

type IMTAiService interface {
	TransferOldFixMT(file multipart.File, userId int64, quota int) (filename string, direction string, err error)
}

type MTAiService struct {
	repository repositories.IUploadRepository
}

func NewMTAiService() IMTAiService {
	return &MTAiService{repository: repositories.NewUploadRepository()}
}

func (m MTAiService) TransferOldFixMT(file multipart.File, userId int64, quota int) (filename string, direction string, err error) {
	fileUrl, filename := SaveFile2Url(file)
	var mediaInfoList []map[string]interface{}
	t := map[string]interface{}{
		"media_data": fileUrl,
		"media_profiles": map[string]interface{}{
			"media_data_type": "url",
		},
	}
	mediaInfoList = append(mediaInfoList, t)
	transferNeedMT := map[string]interface{}{
		"parameter": map[string]interface{}{
			"rsp_media_type": "url",
		},
		"extra":           map[string]interface{}{},
		"media_info_list": mediaInfoList,
	}
	println(transferNeedMT)
	return m.transfer(transferNeedMT, userId, quota)
}

func (d MTAiService) transfer(transferNeedMT map[string]interface{}, userId int64, quota int) (filename string, direction string, err error) {
	b, json_err := json.Marshal(transferNeedMT) //json化结果集
	if json_err != nil {
		fmt.Println("encoding faild")
	} else {
		fmt.Println(string(b))
	}
	println(apiOldFixMT + "?api_key=" + oldFixMTAppKey + "&api_secret=" + oldFixMTAppSecret)
	println(strings.NewReader(string(b)))
	resp, err := http.Post(apiOldFixMT+"?api_key="+oldFixMTAppKey+"&api_secret="+oldFixMTAppSecret,
		"application/json",
		strings.NewReader(string(b)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	//client := &http.Client{}
	//
	//req, err := http.NewRequest(http.MethodPost, apiOldFixMT, strings.NewReader("image="+transferNeedMT.fileUrl))
	//if err != nil {
	//    panic(err)
	//}
	//req.Header.Set("api-key", apiKey)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//    panic(err)
	//}
	//body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(resp)
		fmt.Println(string(body))
		return filename, direction, errors.New("解析出现错误，请重试")
	}
	cc, err := simplejson.NewJson(body)
	fmt.Println(cc)
	if err != nil {
		fmt.Println(resp)
		fmt.Println(string(body))
		return filename, direction, errors.New("json解析出现错误，请重试")
	}

	resUrl := cc.Get("media_info_list").GetIndex(0).Get("media_data").MustString()
	fmt.Println(resUrl)
	SaveImgUrlToLocal(resUrl, filename, Out)

	db := common.NewDbEngine()
	session := db.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}
	//添加数据库记录
	_, err = d.repository.Create(&datamodels.Upload{
		File:   filename,
		UserId: userId,
		Type:   oldFixMT,
		Plate:  datamodels.PlateDeepAi,
	})
	if err != nil {
		session.Rollback()
		return
	}

	//去除额度
	if quota != 0 {
		userService := NewUserService()
		_, err = userService.DescQuotaByUserId(userId)
		if err != nil {
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}

	filename += ImgType
	return
}
