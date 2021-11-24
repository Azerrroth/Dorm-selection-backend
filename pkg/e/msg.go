package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_EXIST_USER:               "已存在该用户名",
	ERROR_NOT_EXIST_USER:           "该用户不存在",
	ERROR_PASSWORD_ERROR:           "用户名或密码错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token验证失败，请重新登录",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	ERROR_AUTH_COOKIE:              "会话不存在或过期，请重新登录",
	ERROR_NO_MORE_ITEMS:            "没有更多的项目",
	ERROR_STUDENT_NOT_EXIST:        "学生不存在",
	ERROR_CERTIFY_CODE_NOT_MATCH:   "验证码不匹配",
	ERROR_USER_IN_BUILDING:         "用户已经有宿舍",
	ERROR_ROOM_NOT_EXIST:           "没有足够数量的宿舍",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
