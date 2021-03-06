package services

import (
	"bufio"
	"comic/common"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
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

func SaveFile2Url(file multipart.File) (fileUrl, filename string) {
	filename = SaveImgFileToLocal(file, In)
	return GetFileUrl(filename, In), filename
}

func SaveImgFileToLocal(file multipart.File, path string) string {
	name := common.GetRandomString(FileNameNum)
	out, err := os.OpenFile(dir+path+name+ImgType, os.O_WRONLY|os.O_CREATE, 06666)
	defer out.Close()
	if err != nil {
		return ""
	}
	io.Copy(out, file)
	return name
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
