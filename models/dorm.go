package models

type Dorm struct {
	Model

	BuildingID int  `gorm:"NOT NULL" json:"building_id"`
	RoomID     int  `json:"room_id"`
	BedID      int  `gorm:"NOT NULL" json:"bed_id"`
	Available  bool `json:"available"` // if available, the bed has no user.
}

func GetDorms(maps interface{}) (dorms []Dorm) {
	dormDB.Where(maps).Find(&dorms)

	return
}

func GetDormCount(maps interface{}) (count int64) {
	dormDB.Model(&Dorm{}).Where(maps).Count(&count)

	return
}

func GetDormPage(skip int, nums int) (dorms []Dorm) {
	dormDB.Offset(skip).Limit(nums).Find(&dorms)

	return
}

func AddDormItem(dorm *Dorm) bool {
	result := dormDB.Create(&dorm)
	return result.RowsAffected != 0
}
