package models

import (
	"go-gin-example/pkg/gender"
)

type Building struct {
	Model

	Name string `json:"name"`
}
type APIBuilding struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetBuildings() ([]Building, error) {
	var buildings []Building
	err := dormDB.Find(&buildings).Error
	return buildings, err
}

func GetBuildingsList() ([]Building, error) {
	var buildings []Building
	err := dormDB.Select("id, name").Find(&buildings).Error
	return buildings, err
}

func AddBuilding(name string) bool {
	building := Building{Name: name}
	result := dormDB.Create(&building)
	return result.RowsAffected > 0
}

func GetAvailableBedsInBuilding(buildingID uint) (maleBeds int, maleTotal int, femaleBeds int, femaleTotal int, err error) {
	var building Building
	err = dormDB.First(&building, buildingID).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var rooms []Room
	err = dormDB.Where("building_id = ?", buildingID).Find(&rooms).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	for _, room := range rooms {
		if room.Gender == gender.FEMALE {
			femaleBeds += int(room.AvailableBeds)
			femaleTotal += int(room.TotalBeds)
		} else {
			maleBeds += int(room.AvailableBeds)
			maleTotal += int(room.TotalBeds)
		}
	}

	return
}

func GetTotalBedsInBuilding(buildingID uint) (maleBeds int, femaleBeds int, err error) {
	var building Building
	err = dormDB.First(&building, buildingID).Error
	if err != nil {
		return 0, 0, err
	}

	var rooms []Room
	err = dormDB.Where("building_id = ?", buildingID).Find(&rooms).Error
	if err != nil {
		return 0, 0, err
	}
	for _, room := range rooms {
		if room.Gender == gender.FEMALE {
			femaleBeds += int(room.TotalBeds)
		} else {
			maleBeds += int(room.TotalBeds)
		}
	}

	return
}
