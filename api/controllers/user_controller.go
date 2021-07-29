package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/common/wechat"
	"comic/datamodels"
	"comic/services"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

type UserController struct {
	Ctx iris.Context
}

var form map[string]string

func (c *UserController) parseForm() {
	err := json.NewDecoder(c.Ctx.Request().Body).Decode(&form)
	if err != nil {
		panic(err.Error())
	}
	defer c.Ctx.Request().Body.Close()
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany(http.MethodPost, "/login", "Login")
}

func (c *UserController) Login() common.Response {
	c.parseForm()
	code := form["code"]

	fmt.Println(code)
	mini := wechat.NewMini()
	res, err := mini.Code2Session(code)
	if err != nil {
		c.Ctx.Application().Logger().Error(err)
		return common.ReErrorMsg(err.Error())
	}

	service := services.NewUserService()

	user := datamodels.User{
		Openid: res.OpenID,
	}
	has := service.Get(&user)
	if !has {
		//创建新用户
		id, err := service.NewUser(&datamodels.User{
			Openid:  res.OpenID,
			UnionId: res.UnionID,
		})
		if err != nil {
			return common.ReErrorMsg(err.Error())
		}
		user = datamodels.User{
			Id: id,
		}
		service.Get(&user)
	}
	fmt.Println(user)

	return common.ReSuccessData(map[string]string{
		"token": middleware.BuildToken(&user),
	})
}
