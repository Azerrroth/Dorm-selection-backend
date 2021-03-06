package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_USER     = 10001
	ERROR_NOT_EXIST_USER = 10002
	ERROR_PASSWORD_ERROR = 10003

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_AUTH_COOKIE              = 20005

	ERROR_NO_MORE_ITEMS          = 30001
	ERROR_STUDENT_NOT_EXIST      = 30002
	ERROR_CERTIFY_CODE_NOT_MATCH = 30003
	ERROR_USER_IN_BUILDING       = 30004
	ERROR_ROOM_NOT_EXIST         = 30005
)
