package v1

import (
	"fmt"
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. Check certify code & student id is matched
// 2. Check users is not in the building
// 3. Check users are same gender
// 4. Make an order
// 5. Check available room for them
//    If exists, response. Else response error
func BookOrder(c *gin.Context) {
	var err error
	json := make(map[string]interface{})
	c.BindJSON(&json)
	code := e.INVALID_PARAMS

	buildingId := uint(json["buildingId"].(float64))
	gender := uint(json["gender"].(float64))
	usersNum := uint(json["usersNum"].(float64))

	users := json["users"].([]interface{})
	// usersInfo := []map[string]interface{}{}
	userInfoList := []models.UserInfo{}

	for _, user := range users {
		tempInfo := make(map[string]interface{})
		tempInfo["userCertifyCode"] = user.(map[string]interface{})["userCertifyCode"]
		tempInfo["userStudentId"] = user.(map[string]interface{})["userStudentId"]
		// fmt.Printf("%v %v\n", tempInfo["userCertifyCode"], tempInfo["userStudentId"])
		// usersInfo = append(usersInfo, tempInfo)

		tempUserInfo, err := models.GetUserInformationByStudentID(tempInfo["userStudentId"].(string))
		if err != nil {
			fmt.Println(err)
			code = e.ERROR_STUDENT_NOT_EXIST
			errorResponse(c, code)
			return
		} else {
			userInfoList = append(userInfoList, tempUserInfo)
		}
		// 1. Check certify code & student id is matched
		if tempUserInfo.CertifyCode != tempInfo["userCertifyCode"].(string) {
			code = e.ERROR_CERTIFY_CODE_NOT_MATCH
			errorResponse(c, code)
			return
		}
	}

	// 2. Check users is not in the building
	for _, userInfo := range userInfoList {
		u2rs, _ := models.GetValidUser2RoomByUserId(userInfo.ID)
		// fmt.Println(len(u2rs))
		if len(u2rs) != 0 {
			code = e.ERROR_USER_IN_BUILDING
			errorResponse(c, code)
			return
		}
		// 3. Check users are same gender
		if userInfo.Gender != gender {
			code = e.ERROR_USER_IN_BUILDING
			errorResponse(c, code)
			return
		}
	}

	// 4. Make an order
	order := models.Order{
		UserId:        userInfoList[0].ID,
		BuildingId:    buildingId,
		Gender:        gender,
		NumberOfGroup: usersNum,
	}
	order.ID, err = models.AddOrder(&order)
	if err != nil {
		code = e.ERROR
		errorResponse(c, code)
		return
	}
	for _, user := range userInfoList {
		orderDetail := models.OrderDetail{
			ResidentId: user.ID,
			OrderId:    order.ID,
		}
		models.AddOrderDetail(&orderDetail)
	}
	// 5. Check available room for them
	//    If exists, response. Else response error

	rooms, _ := models.GetRoomWithAvailableBedsByBuildingIdAndGender(buildingId, usersNum, gender)
	if len(rooms) == 0 {
		order.IsSuccess = false
		models.UpdateOrder(&order)
		code = e.ERROR_ROOM_NOT_EXIST
		errorResponse(c, code)
		return
	}
	choice := rand.Intn(len(rooms))
	room := rooms[choice]
	room.AvailableBeds -= usersNum
	models.UpdateRoom(room)

	order.IsSuccess = true
	models.UpdateOrder(&order)

	for _, user := range userInfoList {
		models.AddUser2Room(user.ID, room.ID)
	}
	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func errorResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": nil,
	})
}
