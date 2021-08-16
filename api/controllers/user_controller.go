package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/common/wechat"
	"comic/datamodels"
	"comic/repositories"
	"comic/services"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) parseForm() {
	err := json.NewDecoder(c.Ctx.Request().Body).Decode(&form)
	if err != nil {
		panic(err.Error())
	}
	defer c.Ctx.Request().Body.Close()
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany(http.MethodPost, "/login", "Login")
	b.HandleMany(http.MethodGet, "/info", "UserInfo", middleware.AuthTokenHandler().Serve)
}

func (c *UserController) Login() common.Response {
	appId, err := strconv.ParseInt(c.Ctx.GetHeader("app_id"), 10, 64)
	if err != nil {
		c.Ctx.Application().Logger().Error(err)
		return common.ReErrorMsg(err.Error())
	}

	c.parseForm()
	code := form["code"]

	mini := wechat.NewMini(appId)
	res, err := mini.Code2Session(code)
	if err != nil {
		c.Ctx.Application().Logger().Error(err)
		return common.ReErrorMsg(err.Error())
	}

	service := services.NewUserService()

	user := datamodels.User{
		Openid: res.OpenID,
		AppId:  appId,
	}
	has := service.Get(&user)
	if !has {

		AppRepository := repositories.NewAppRepository()
		has, app := AppRepository.Get(appId)
		if !has {
			return common.ReErrorMsg("App不存在:" + strconv.FormatInt(appId, 10))
		}

		//创建新用户
		id, err := service.NewUser(&datamodels.User{
			Openid:  res.OpenID,
			UnionId: res.UnionID,
			AppId:   appId,
			Quota:   app.DefaultQuota,
		})
		if err != nil {
			return common.ReErrorMsg(err.Error())
		}
		user = datamodels.User{
			Id: id,
		}
		service.Get(&user)
	}

	return common.ReSuccessData(map[string]string{
		"token": middleware.BuildToken(&user),
	})
}

func (c *UserController) UserInfo() common.Response {
	user := middleware.ParseTokenToUser(c.Ctx)
	userService := services.NewUserService()
	userService.Get(&datamodels.User{
		Id: user.Id,
	})

	return common.ReSuccessData(map[string]int64{
		"id":    user.Id,
		"quota": user.Quota,
	})
}
