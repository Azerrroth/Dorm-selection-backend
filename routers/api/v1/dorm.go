package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"go-gin-example/models"
	"go-gin-example/pkg/e"
)

func GetBuildingList(c *gin.Context) {
	code := e.INVALID_PARAMS
	buildings, err := models.GetBuildings()
	if err == nil {
		code = e.SUCCESS
	} else {
		code = e.ERROR
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": buildings,
	})
}

func GetBuildingStatus(c *gin.Context) {
	code := e.INVALID_PARAMS
	data := make(map[string]interface{})

	building_id := c.Query("building_id")
	valid := validation.Validation{}
	valid.Required(building_id, "building_id").Message("楼栋id不能为空")
	buildingId, _ := strconv.Atoi(building_id)

	if !valid.HasErrors() {
		maleAva, maleTotal, femaleAva, femaleTotal, err := models.GetAvailableBedsInBuilding(uint(buildingId))
		// fmt.Println(maleAva, maleTotal, femaleAva, femaleTotal)
		data["building_id"] = building_id
		data["male_available"] = maleAva
		data["male_total"] = maleTotal
		data["female_available"] = femaleAva
		data["female_total"] = femaleTotal
		if err == nil {
			code = e.SUCCESS
		} else {
			code = e.ERROR
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func GetBuildingsStatus(c *gin.Context) {
	code := e.INVALID_PARAMS
	result := []map[string]interface{}{}
	buildings, err := models.GetBuildings()
	for _, building := range buildings {
		data := make(map[string]interface{})
		maleAva, maleTotal, femaleAva, femaleTotal, _ := models.GetAvailableBedsInBuilding(building.ID)
		data["building_id"] = building.ID
		data["building_name"] = building.Name
		data["male_available"] = maleAva
		data["male_total"] = maleTotal
		data["female_available"] = femaleAva
		data["female_total"] = femaleTotal
		result = append(result, data)
	}

	if err == nil {
		code = e.SUCCESS
	} else {
		code = e.ERROR
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": result,
	})
}

func GetUser2RoomInfo(c *gin.Context) {
	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	user_dorm_info := make(map[string]interface{})

	user_id := c.Query("user_id")
	if user_id == "" {
		user_id = c.GetHeader("x-user-id")
	}

	valid := validation.Validation{}
	valid.Required(user_id, "user_id").Message("用户id不能为空")
	userId, _ := strconv.Atoi(user_id)
	u2rs, _ := models.GetValidUser2RoomByUserId(uint(userId))
	if len(u2rs) != 0 {
		code = e.SUCCESS
		room := models.GetRoomByID(u2rs[0].RoomId)
		user_dorm_info["roomId"] = room.ID
		user_dorm_info["roomName"] = room.Name
		user_dorm_info["roommatesNum"] = room.TotalBeds - room.AvailableBeds - room.InvalidBeds
		data["userDormInfo"] = user_dorm_info
	} else {
		code = e.SUCCESS
	}
	data["userDormInfo"] = user_dorm_info
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func CheckOutRoom(c *gin.Context) {
	code := e.INVALID_PARAMS
	json := make(map[string]interface{})
	valid := validation.Validation{}
	err := c.BindJSON(&json)
	var user_id, room_id string
	if err == nil {
		user_id = c.GetHeader("x-user-id")
		room_id = json["roomId"].(string)
		if user_id == "" {
			user_id = json["userId"].(string)
		}
		valid.Required(user_id, "user_id").Message("用户id不能为空")
		valid.Required(room_id, "room_id").Message("宿舍id不能为空")
	}

	if !valid.HasErrors() {
		userId, _ := strconv.Atoi(user_id)
		roomId, _ := strconv.Atoi(room_id)
		u2rs, _ := models.GetValidUser2RoomByUserId(uint(userId))
		if len(u2rs) != 0 {
			for _, u2r := range u2rs {
				if u2r.RoomId == uint(roomId) {
					u2r.IsValid = false
					models.UpdateUser2Room(&u2r)
					room := models.GetRoomByID(uint(roomId))
					room.AvailableBeds += 1
					models.UpdateRoom(&room)
					code = e.SUCCESS
					break
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}
