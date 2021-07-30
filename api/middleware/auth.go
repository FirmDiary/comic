package middleware

import (
	"comic/common"
	"comic/datamodels"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"time"
)

const (
	Ttl = 10080 //有效期 7天
)

func AuthTokenHandler() *jwt.Middleware {
	// j2 对比 j 添加了错误处理函数
	return jwt.New(jwt.Config{
		// 注意，新增了一个错误处理函数
		ErrorHandler: func(ctx iris.Context, err error) {
			if err == nil {
				return
			}
			ctx.StopExecution()
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(common.Response{
				Code: 401,
				Msg:  "auth failed",
				Data: nil,
			})
		},
		// 设置一个函数返回秘钥，关键在于return []byte("这里设置秘钥")
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},

		// 设置一个加密方法
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func BuildToken(user *datamodels.User) (tokenString string) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 根据需求，可以存一些必要的数据
		"id":     user.Id,
		"openid": user.Openid,
		"app_id": user.AppId,
		// 签发时间
		"iat": time.Now().Unix(),
		// 设定过期时间，便于测试，设置1分钟过期
		"exp": time.Now().Add(1 * time.Minute * time.Duration(Ttl)).Unix(),
	})

	// 使用设置的秘钥，签名生成jwt字符串
	tokenString, _ = token.SignedString([]byte("My Secret"))
	return
}

func print_json(m map[string]interface{}) {
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float", int64(vv))
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case nil:
			fmt.Println(k, "is nil", "null")
		case map[string]interface{}:
			fmt.Println(k, "is an map:")
			print_json(vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle ", fmt.Sprintf("%T", v))
		}
	}
}

func ParseTokenToUser(ctx iris.Context) *datamodels.User {
	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)

	return &datamodels.User{
		Id:     int64(jwtInfo["id"].(float64)),
		Openid: jwtInfo["openid"].(string),
		AppId:  int64(jwtInfo["app_id"].(float64)),
	}
}
