package models

import "fmt"

type OrderDetail struct {
	Model

	ResidentId uint `json:"resident_id"`
	OrderId    uint `json:"order_id"`
}

func GetAllResidentsByOrderId(orderId string) ([]*OrderDetail, error) {
	var orderDetails []*OrderDetail
	err := orderDB.Where("order_id = ?", orderId).Find(&orderDetails).Error
	return orderDetails, err
}

func GetResidentsInfoByOrderId(orderId uint) ([]UserInfo, error) {
	userTableName := tablePrefix + "users"
	orderDetailTableName := tablePrefix + "order_details"

	sql := fmt.Sprintf(`SELECT a.* FROM 
	(SELECT id, student_id, name, gender, mail, authority FROM %s)
	as a right join
	(SELECT resident_id FROM %s WHERE order_id = (?)) as b
	on a.id = b.resident_id`, userTableName, orderDetailTableName)
	var result []UserInfo

	err := userDB.Raw(sql, orderId).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AddOrderDetail(orderDetail *OrderDetail) error {
	return orderDB.Create(orderDetail).Error
}
