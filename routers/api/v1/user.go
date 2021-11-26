package v1

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/validation"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
)

var jwtValidityPeriod = setting.JwtValidityPeriod
var LoginList map[string]interface{}

func GetUser(c *gin.Context) {
	username := c.Query("username")

	maps := make(map[string]interface{})
	// data := make(map[string]interface{})

	if username != "" {
		maps["username"] = username
	}

}

func Register(c *gin.Context) {
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	// 检查表单
	username := c.Query("username")
	password := c.Query("password")

	if username == "" {
		if err == nil {
			username = json["username"].(string)
		} else {
			username = c.PostForm("username")
		}
	}

	if password == "" {
		if err == nil {
			password = json["password"].(string)
		} else {
			password = c.PostForm("password")
		}
	}
	valid := validation.Validation{}
	valid.Required(username, "username").Message("用户名不能为空")
	valid.MaxSize(username, 30, "username").Message("用户名最长为30字符")
	valid.Required(password, "password").Message("密码不能为空")
	valid.MaxSize(password, 30, "password").Message("密码最长为30字符")
	valid.MinSize(password, 6, "password").Message("密码最少为6字符")

	crypted, salt, err := util.Encrypt(password)
	if err != nil {
		log.Fatal(2, "Fail to get encrypted password: %v", err)
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistUserByUsername(username) {
			code = e.SUCCESS
			models.AddUser(username, string(crypted), string(salt))
		} else {
			code = e.ERROR_EXIST_USER
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

func Login(c *gin.Context) {
	json := make(map[string]interface{})
	err := c.ShouldBind(&json)
	// 检查表单
	username := c.Query("username")
	password := c.Query("password")

	if username == "" {
		if err == nil {
			username = json["username"].(string)
		} else {
			username = c.PostForm("username")
		}
	}

	if password == "" {
		if err == nil {
			password = json["password"].(string)
		} else {
			password = c.PostForm("password")
		}
	}

	valid := validation.Validation{}
	valid.Required(username, "username").Message("用户名不能为空")
	valid.Required(password, "password").Message("密码不能为空")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	// session := sessions.Default(c)
	if !valid.HasErrors() {
		if !models.ExistUserByUsername(username) {
			code = e.ERROR_NOT_EXIST_USER
		} else {
			user, right := models.ExistUserByUsernameAndPassword(username, password)
			userInfo, _ := models.GetUserInformationByUsername(username)
			if right {
				// 登录成功
				code = e.SUCCESS
				token, err := util.GenerateToken(username, password, jwtValidityPeriod)
				// session.Set("uid", username)
				// session.Set("status", "online")
				// session.Save()
				if err != nil {
					code = e.ERROR_AUTH_TOKEN
				} else {
					c.Header("new-token", token)
					data["token"] = token
					data["userInfo"] = userInfo
					currentTime := time.Now()
					user.LoginAt = &currentTime
					models.UpdateUser(&user)
					// data["uid"] = session.Get("uid")
				}
				// Open session
			} else {
				// 密码错误
				code = e.ERROR_PASSWORD_ERROR
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func UpdateCertifyCode(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetHeader("x-user-id"))
	valid := validation.Validation{}
	valid.Required(id, "id").Message("未携带用户ID头")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if !valid.HasErrors() {
		if !models.ExistUserByID(uint(id)) {
			code = e.ERROR_NOT_EXIST_USER
		} else {
			code = e.SUCCESS
			data["certifyCode"], _ = models.UpdateCertifyCode(uint(id))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func UpdateUserProfile(c *gin.Context) {
	code := e.INVALID_PARAMS
	json := make(map[string]interface{})
	valid := validation.Validation{}
	err := c.BindJSON(&json)
	var user_id string

	student_id := json["student_id"].(string)
	gender := json["gender"].(float64)
	name_str := json["name"].(string)
	mail := json["mail"].(string)

	if err == nil {
		user_id = c.GetHeader("x-user-id")
		if user_id == "" {
			user_id = json["userId"].(string)
		}
		valid.Required(user_id, "user_id").Message("用户id不能为空")
	}

	if !valid.HasErrors() {
		id, _ := strconv.Atoi(user_id)

		if !models.ExistUserByID(uint(id)) {
			code = e.ERROR_NOT_EXIST_USER
		} else {
			user := models.GetUserByID(uint(id))
			user.Name = &name_str
			user.Gender = uint(gender)
			user.Mail = &mail
			user.StudentId = &student_id

			models.UpdateUser(&user)
			code = e.SUCCESS
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": json,
	})
}
