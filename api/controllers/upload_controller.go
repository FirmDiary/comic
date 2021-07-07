package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type UploadController struct {
	Ctx iris.Context
}

func (u *UploadController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany("POST", "/transfer", "Transfer", middleware.AuthTokenHandler().Serve)
}

func (u *UploadController) Transfer() common.Response {
	transferTypes := u.Ctx.GetHeader("transferTypes")
	if transferTypes == "" {
		transferTypes = "1"
	}
	transferType, err := strconv.Atoi(transferTypes)
	if err != nil {
		return common.ReErrorMsg(err.Error())
	}

	file, _, err := u.Ctx.FormFile("file")
	defer file.Close()

	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return common.ReErrorMsg(err.Error())
	}

	user := middleware.ParseTokenToUser(u.Ctx)

	service := services.NewUploadService()
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
