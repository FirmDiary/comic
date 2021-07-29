package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/services/transfer"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"mime/multipart"
	"net/http"
	"strconv"
)

type UploadController struct {
	Ctx iris.Context
}

func (u *UploadController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany(http.MethodPost, "/transferU2", "TransferU2", middleware.AuthTokenHandler().Serve)
	b.HandleMany(http.MethodPost, "/transferOldFix", "TransferOldFix", middleware.AuthTokenHandler().Serve)
}

func getFile(u *UploadController) (multipart.File, error) {
	file, _, err := u.Ctx.FormFile("file")
	defer file.Close()

	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return file, err
	}
	return file, err
}

func (u *UploadController) TransferOldFix() common.Response {
	file, err := getFile(u)
	if err != nil {
		return common.ReErrorMsg(err.Error())
	}

	user := middleware.ParseTokenToUser(u.Ctx)

	service := transfer.NewDeepAiService()
	filename, err := service.TransferOldFix(file, user.Id)

	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return common.ReErrorMsg(err.Error())
	}

	return common.ReSuccessData(map[string]string{
		"filename": filename,
	})
}

func (u *UploadController) TransferU2() common.Response {
	transferTypes := u.Ctx.GetHeader("transfer_type")
	fmt.Println(transferTypes)
	if transferTypes == "" {
		transferTypes = "1"
	}
	transferType, err := strconv.Atoi(transferTypes)
	if err != nil {
		return common.ReErrorMsg(err.Error())
	}

	file, err := getFile(u)
	if err != nil {
		return common.ReErrorMsg(err.Error())
	}

	user := middleware.ParseTokenToUser(u.Ctx)

	service := transfer.NewUploadService()
	path, err := service.Transfer(file, user.Id, transferType)

	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return common.ReErrorMsg(err.Error())
	}

	return common.ReSuccessData(map[string]string{
		"path": path,
	})
}
