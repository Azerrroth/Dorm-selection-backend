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
	orderDB.Where(maps).Find(&orders)
	return
}

func GetOrderCount(maps interface{}) (count int64) {
	orderDB.Model(&Order{}).Where(maps).Count(&count) // 查询总数
	return
}

func AddOrder(order *Order) (uint, error) {
	result := orderDB.Create(order)
	return order.ID, result.Error
}

func UpdateOrder(order *Order) error {
	err := orderDB.Save(order).Error
	return err
}
