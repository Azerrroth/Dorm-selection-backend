package models

type Room struct {
	Model

	BuildingId    int    `json:"building_id"`
	Name          string `json:"name"`
	Gender        uint   `json:"gender"` // 1 男，2 女
	TotalBeds     uint   `json:"total_beds"`
	AvailableBeds uint   `json:"available_beds"`
	InvalidBeds   uint   `gorm:"default:0" json:"invalid_beds"`
}

func GetRoomCountByBuildingId(buildingId int) (count int64) {
	dormDB.Model(&Room{}).Where("building_id = ?", buildingId).Count(&count)
	return
}

func GetRooms(maps interface{}) (rooms []*Room) {
	dormDB.Find(&rooms)
	return
}

func GetRoomByID(id uint) (room Room) {
	dormDB.Where("id = ?", id).First(&room)
	return
}

func GetRoomCount(maps interface{}) (count int64) {
	dormDB.Model(&Room{}).Where(maps).Count(&count)
	return
}

func GetRoomWithAvailableBedsByBuildingId(buildingId uint, userNum uint) (rooms []*Room, err error) {
	err = dormDB.Where("building_id = ? AND available_beds > ?", buildingId, userNum).Find(&rooms).Error
	return
}

func GetRoomWithAvailableBedsByBuildingIdAndGender(buildingId uint, userNum uint, gender uint) (rooms []*Room, err error) {
	err = dormDB.Where("building_id = ? AND available_beds > ? AND gender = ?", buildingId, userNum, gender).Find(&rooms).Error
	return
}

func AddRoom(room *Room) (err error) {
	err = dormDB.Create(room).Error
	return
}

func UpdateRoom(room *Room) (err error) {
	err = dormDB.Save(room).Error
	return
}
