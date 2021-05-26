package common

import (
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ReSuccessData(data interface{}) Response {
	return Response{0, "", data}
}

func ReSuccessMsg(msg string) Response {
	return Response{1, msg, ""}
}

func ReErrorMsg(errMsg string) Response {
	return Response{http.StatusInternalServerError, errMsg, ""}
}
