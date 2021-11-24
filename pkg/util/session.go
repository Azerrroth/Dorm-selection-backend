package util

import (
	"github.com/gin-gonic/gin"

	"go-gin-example/pkg/setting"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var sessionKey = []byte(setting.JwtSecret)
var Store = cookie.NewStore(sessionKey)

func Get(c *gin.Context, key interface{}) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

func Set(c *gin.Context, key interface{}, val interface{}) {
	session := sessions.Default(c)
	session.Set(key, val)
}

func Delete(c *gin.Context, key interface{}) {
	session := sessions.Default(c)
	session.Delete(key)
}

func Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func Save(c *gin.Context) error {
	session := sessions.Default(c)
	return session.Save()
}
