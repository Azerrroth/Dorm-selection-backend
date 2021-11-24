package session

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-gin-example/pkg/e"
)

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {

		var code int
		var data interface{}
		session := sessions.Default(c)
		code = e.SUCCESS
		fmt.Println(session.Get("uid"))
		fmt.Println(session.Get("status"))
		if cookie, err := c.Request.Cookie("session"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "" {
				code = e.ERROR_AUTH_COOKIE
			}
		} else {
			code = e.ERROR_AUTH_COOKIE
		}

		// if session.Get("uid") == nil {
		// 	code = e.INVALID_PARAMS
		// }
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
