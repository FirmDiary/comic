package controllers

import (
	"comic/api/middleware"
	"comic/common"
	"comic/datamodels"
	"comic/repositories"
	"comic/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
	"strconv"
)

type InviteController struct {
	Ctx iris.Context
}

func (c *InviteController) BeforeActivation(b mvc.BeforeActivation) {
	b.HandleMany(http.MethodGet, "/success", "InviteSuccess", middleware.AuthTokenHandler().Serve)
}

/**
完成邀请
*/
func (c *InviteController) InviteSuccess() common.Response {
	inviteId, err := strconv.ParseInt(c.Ctx.FormValue("invite_id"), 10, 64)

	if err != nil {
		c.Ctx.Application().Logger().Error(err)
		return common.ReErrorMsg(err.Error())
	}

	user := middleware.ParseTokenToUser(c.Ctx)

	if user.Id == inviteId {
		return common.ReSuccessData("自己邀请自己")
	}

	InviteRepository := repositories.NewInviteLogRepository()
	log := &datamodels.InviteLog{
		UserId:   user.Id,
		InviteId: inviteId,
	}
	has := InviteRepository.Get(log)

	if has {
		return common.ReSuccessData("已经邀请过")
	}

	InviteRepository.Add(log)

	userService := services.NewUserService()
	userService.AddQuotaByUserId(inviteId, 2)

	return common.ReSuccessData("")
}
