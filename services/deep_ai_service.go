package services

import (
	"comic/common"
	"comic/datamodels"
	"comic/repositories"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	oldFix  = 1 //老照片修复
	waifu2x = 2 //通过图片像素放大两倍使其变清晰

	apiOldFix = "https://api.deepai.org/api/colorizer"
	api2x     = "https://api.deepai.org/api/waifu2x"

	apiKey = "764f9fc0-97de-4fe3-bf26-0c4ee9052139"
)

type IDeepAiService interface {
	TransferOldFix(file multipart.File, userId int64, quota int) (filename string, direction string, err error)
	Transfer2x(fileUrl string, userId int64, quota int) (filename string, direction string, err error)
}

type DeepAiService struct {
	repository repositories.IUploadRepository
}

type TransferNeed struct {
	fileUrl      string
	userId       int64
	transferType int
	api          string
	quota        int
	direction    bool
	filename     string
}

func NewDeepAiService() IDeepAiService {
	return &DeepAiService{repository: repositories.NewUploadRepository()}
}

func (d DeepAiService) TransferOldFix(file multipart.File, userId int64, quota int) (filename string, direction string, err error) {
	fileUrl, filename := saveFile2Url(file)
	transferNeed := NewTransferNeed(fileUrl, userId, oldFix, filename, quota)

	image1, _ := mergi.Import(impexp.NewFileImporter(dir + Out + filename + ImgType))
	image2, _ := mergi.Import(impexp.NewFileImporter(dir + In + filename + ImgType))

	horizontalImage, _ := mergi.Merge("TT", []image.Image{image1, image2})
	err = mergi.Export(impexp.NewFileExporter(horizontalImage, dir+In+"666.png"))
	fmt.Println(666)
	fmt.Println(err)

	verticalImage, _ := mergi.Merge("TB", []image.Image{image1, image2})
	err = mergi.Export(impexp.NewFileExporter(verticalImage, dir+In+"777.png"))
	fmt.Println(777)
	fmt.Println(err)

	return d.transfer(transferNeed)
}

func (d DeepAiService) Transfer2x(fileUrl string, userId int64, quota int) (filename string, direction string, err error) {
	index := strings.LastIndex(fileUrl, ".")
	if index == -1 {
		fmt.Println("获取图片名称出现错误" + fileUrl)
		err = errors.New("获取图片名称出现错误" + fileUrl)
		return
	}
	filename = fileUrl[index-FileNameNum : index]
	transferNeed := NewTransferNeed(fileUrl, userId, waifu2x, filename, quota)
	return d.transfer(transferNeed)
}

func NewTransferNeed(fileUrl string, userId int64, transferType int, filename string, quota int) (transferNeed *TransferNeed) {
	transferNeed = &TransferNeed{
		fileUrl:      fileUrl,
		userId:       userId,
		transferType: transferType,
		api:          "",
		quota:        quota,
		filename:     filename,
	}
	switch transferNeed.transferType {
	case oldFix:
		transferNeed.api = apiOldFix
		transferNeed.direction = true
	case waifu2x:
		transferNeed.api = api2x
		transferNeed.direction = false
	}
	return
}

func saveFile2Url(file multipart.File) (fileUrl, filename string) {
	filename = SaveImgFileToLocal(file, In)
	return GetFileUrl(filename, In), filename
}

func (d DeepAiService) transfer(transferNeed *TransferNeed) (filename string, direction string, err error) {
	filename = transferNeed.filename

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, transferNeed.api, strings.NewReader("image="+transferNeed.fileUrl))
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
	SaveImgUrlToLocal(resUrl, filename, Out)

	if transferNeed.direction {
		direction = GetImgDirection(transferNeed.fileUrl)
	}

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
		UserId: transferNeed.userId,
		Type:   transferNeed.transferType,
		Plate:  datamodels.PlateDeepAi,
	})
	if err != nil {
		session.Rollback()
		return
	}

	//去除额度
	if transferNeed.quota != 0 {
		userService := NewUserService()
		_, err = userService.DescQuotaByUserId(transferNeed.userId)
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
