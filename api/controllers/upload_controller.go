package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/datamodels"
	"comic/services"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"mime/multipart"
	"net/http"
)

type UploadController struct {
	Ctx iris.Context
}

func (u *UploadController) parseForm() {
	err := json.NewDecoder(u.Ctx.Request().Body).Decode(&form)
	if err != nil {
		panic(err.Error())
	}
	defer u.Ctx.Request().Body.Close()
}

func (u *UploadController) BeforeActivation(b mvc.BeforeActivation) {
	//b.HandleMany(http.MethodPost, "/transferU2", "TransferU2", middleware.AuthTokenHandler().Serve)
	b.HandleMany(http.MethodPost, "/transferOldFix", "TransferOldFix", middleware.AuthTokenHandler().Serve)
	b.HandleMany(http.MethodPost, "/transfer2x", "TransferFileUrl2x", middleware.AuthTokenHandler().Serve)
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

func (u *UploadController) prepare() (file multipart.File, user *datamodels.User, err error) {
	file, err = getFile(u)
	if err != nil {
		return
	}
	user = middleware.ParseTokenToUser(u.Ctx)
	userService := services.NewUserService()
	userService.Get(user)
	return
}

func (u *UploadController) dealErr(err error) common.Response {
	u.Ctx.StatusCode(iris.StatusInternalServerError)
	u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
	return common.ReErrorMsg(err.Error())
}

func (u *UploadController) TransferOldFix() common.Response {
	file, user, err := u.prepare()

	if err != nil {
		return common.ReErrorMsg(err.Error())
	}
	if user.Quota == 0 {
		//额度不足
		return common.ReSuccessData(map[string]int64{
			"quota": -1,
		})
	}

	service := services.NewDeepAiService()
	filename, direction, err := service.TransferOldFix(file, user.Id, 1)

	if err != nil {
		return u.dealErr(err)
	}

	type value interface{}
	return common.ReSuccessData(map[string]value{
		"filename":  filename,
		"direction": direction,
	})
}

func (u *UploadController) TransferFileUrl2x() common.Response {
	user := middleware.ParseTokenToUser(u.Ctx)
	userService := services.NewUserService()
	userService.Get(user)

	u.parseForm()
	fileUrl := form["url"]
	useQuota := form["use_quota"]

	quota := 0
	if useQuota == "1" {
		quota = 1
		if user.Quota == 0 {
			//额度不足
			return common.ReSuccessData(map[string]int64{
				"quota": -1,
			})
		}
	}

	service := services.NewDeepAiService()
	filename, _, err := service.Transfer2x(fileUrl, user.Id, quota)

	if err != nil {
		return u.dealErr(err)
	}

	type value interface{}
	return common.ReSuccessData(map[string]value{
		"filename": filename,
	})
}

//
//func (u *UploadController) TransferU2() common.Response {
//    transferTypes := u.Ctx.GetHeader("transfer_type")
//    fmt.Println(transferTypes)
//    if transferTypes == "" {
//        transferTypes = "1"
//    }
//    transferType, err := strconv.Atoi(transferTypes)
//    if err != nil {
//        return common.ReErrorMsg(err.Error())
//    }
//
//    file, err := getFile(u)
//    if err != nil {
//        return common.ReErrorMsg(err.Error())
//    }
//
//    user := middleware.ParseTokenToUser(u.Ctx)
//
//    service := services.NewUploadService()
//    path, err := service.Transfer(file, user.Id, transferType)
//
//    if err != nil {
//        u.Ctx.StatusCode(iris.StatusInternalServerError)
//        u.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
//        return common.ReErrorMsg(err.Error())
//    }
//
//    return common.ReSuccessData(map[string]string{
//        "path": path,
//    })
//}
