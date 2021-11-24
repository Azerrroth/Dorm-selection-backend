package models

import (
	"fmt"
	"time"
)

type User2Room struct {
	Model

	UserId      uint       `json:"user_id"`
	RoomId      uint       `json:"room_id"`
	IsValid     bool       `json:"is_valid"`
	ExpiredTime *time.Time `json:"expired_time"`
}

func (u2r *User2Room) TableName() string {
	return tablePrefix + "user2room"
}

func GetUser2RoomByUserId(userId uint) ([]User2Room, error) {
	var u2rs []User2Room
	err := dormDB.Where("user_id = ?", userId).Find(&u2rs).Error
	return u2rs, err
}

func GetUser2RoomByRoomId(roomId uint, expired_time time.Time) ([]User2Room, error) {
	var u2rs []User2Room
	err := dormDB.Where("room_id = ?", roomId).Find(&u2rs).Error
	return u2rs, err
}

func GetValidUser2RoomByUserId(userId uint) ([]User2Room, error) {
	var u2rs []User2Room
	err := dormDB.Where("user_id = ? AND is_valid = ?", userId, true).Find(&u2rs).Error
	return u2rs, err
}

func GetValidUser2RoomByRoomId(roomId uint) ([]User2Room, error) {
	var u2rs []User2Room
	err := dormDB.Where("room_id = ? AND is_valid = ?", roomId, true).Find(&u2rs).Error
	return u2rs, err
}

func GetValidUserInfoByRoomId(roomId uint) ([]UserInfo, error) {
	var result []UserInfo
	u2r := new(User2Room)
	userTableName := tablePrefix + "users"
	u2rTableName := tablePrefix + u2r.TableName()

	sql := fmt.Sprintf(`SELECT a.* FROM 
	(SELECT id, student_id, name, gender, mail, authority FROM %s)
	as a right join
	(SELECT user_id FROM %s WHERE room_id = (?) AND is_valid = TRUE) as b
	on a.id = b.user_id`, userTableName, u2rTableName)
	err := dormDB.Raw(sql, roomId).Scan(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AddUser2Room(userId, roomId uint) error {
	u2r := User2Room{
		UserId:      userId,
		RoomId:      roomId,
		IsValid:     true,
		ExpiredTime: nil,
	}
	return dormDB.Create(&u2r).Error
}

func (u2r *User2Room) AddU2R() error {
	return dormDB.Create(u2r).Error
}

func UpdateUser2Room(u2r *User2Room) error {
	return dormDB.Save(u2r).Error
}
