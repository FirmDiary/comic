package services

import (
	"bufio"
	"comic/common"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var dir, _ = os.Getwd()

const (
	In      = "/upload/in/"  //文件输入目录
	Out     = "/upload/out/" //文件输出目录
	ImgType = ".png"

	DirectionColumn = "column"
	DirectionRow    = "row"

	ImgUrlPrefix = "https://comic-img.zwww.cool" //图片域名

	FileNameNum = 12
)

func GetFileUrl(name string, path string) string {
	var pathMid string
	if path == In {
		pathMid = "/in/"
	} else {
		pathMid = "/out/"
	}
	return ImgUrlPrefix + pathMid + name + ImgType
}

func SaveFile2Url(file multipart.File) (fileUrl, filename, imgBase64 string) {
	filename, imgBase64 = SaveImgFileToLocal(file, In)
	return GetFileUrl(filename, In), filename, imgBase64
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func SaveImgFileToLocal(file multipart.File, path string) (name, imgBase64 string) {
	name = common.GetRandomString(FileNameNum)
	out, err := os.OpenFile(dir+path+name+ImgType, os.O_WRONLY|os.O_CREATE, 06666)
	defer out.Close()
	if err != nil {
		return
	}
	io.Copy(out, file)

	bytes, err := ioutil.ReadFile(dir + path + name + ImgType)
	if err != nil {
		log.Fatal(err)
	}
	imgBase64 += "data:image/png;base64,"
	imgBase64 += toBase64(bytes)

	return name, imgBase64
}

func SaveImgUrlToLocal(fileUrl string, name string, path string) string {
	res, err := http.Get(fileUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	reader := bufio.NewReaderSize(res.Body, 32*1024)

	filename := dir + path + name + ImgType
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	io.Copy(writer, reader)

	return GetFileUrl(name, Out)
}

func GetImgDirection(fileUrl string) (direction string) {
	res, err := http.Get(fileUrl)
	if err != nil {
		fmt.Println("获取图片方向出现错误")
		fmt.Println(err.Error())
		return DirectionColumn
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println("获取图片方向出现错误")
		fmt.Println(err.Error())
		return DirectionColumn
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width > height {
		direction = DirectionColumn
	} else {
		direction = DirectionRow
	}
	return
}

//删除识别的文件
func DelUploadImg(name string) error {
	fileIn := dir + In + name + ImgType
	fileOut := dir + Out + name + ImgType

	var err error
	err = os.Remove(fileIn)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = os.Remove(fileOut)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}
