package models

type Order struct {
	Model

	UserId        uint `json:"user_id"`
	BuildingId    uint `json:"building_id"`
	Gender        uint `json:"gender"`
	NumberOfGroup uint `json:"number_of_group"`
	IsSuccess     bool `json:"is_success"`
}

func GetOrders(maps interface{}) (orders []Order) {
	dormDB.Where(maps).Find(&orders)
	return
}

func GetOrderCount(maps interface{}) (count int64) {
	dormDB.Model(&Order{}).Where(maps).Count(&count) // 查询总数
	return
}

func AddOrder(order *Order) (uint, error) {
	result := dormDB.Create(order)
	return order.ID, result.Error
}

func UpdateOrder(order *Order) error {
	err := dormDB.Save(order).Error
	return err
}
