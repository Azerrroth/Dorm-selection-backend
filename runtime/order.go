package runtime

import (
	"encoding/json"
	"fmt"
	"go-gin-example/models"
	genderConfig "go-gin-example/pkg/gender"
	"log"
	"math/rand"

	"github.com/streadway/amqp"
)

func DeclareQueue(name string) amqp.Queue {
	qu, err := channel.QueueDeclare(name, false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	return qu
}

func PublishMessage(queue amqp.Queue, message string) {
	err := channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(message),
	})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", message)
}

func ConsumeMessage(queue amqp.Queue) {
	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] Received %s", d.Body)
			DealOrder(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func DealOrder(message string) {
	log.Printf(" [x] Received %s", message)

	var err error
	formatMessage := make(map[string]interface{})
	err = json.Unmarshal([]byte(message), &formatMessage)
	failOnError(err, "Failed to parse message")

	buildingId := uint(formatMessage["buildingId"].(float64))
	gender := uint(formatMessage["gender"].(float64))
	usersNum := uint(formatMessage["usersNum"].(float64))
	users := formatMessage["users"].([]interface{})
	userInfoList := []models.UserInfo{}

	for _, user := range users {
		tempInfo := make(map[string]interface{})
		tempInfo["userCertifyCode"] = user.(map[string]interface{})["userCertifyCode"]
		tempInfo["userStudentId"] = user.(map[string]interface{})["userStudentId"]

		tempUserInfo, err := models.GetUserInformationByStudentID(tempInfo["userStudentId"].(string))
		if err != nil {
			fmt.Println(err)
			return
		} else {
			userInfoList = append(userInfoList, tempUserInfo)
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
		return
	}
	choice := rand.Intn(len(rooms))
	room := rooms[choice]
	room.AvailableBeds -= usersNum
	models.UpdateRoom(room)
	// Deal redis data.
	models.MinusBuildingStatus(buildingId, -int(usersNum), gender == genderConfig.MALE)

	order.IsSuccess = true
	models.UpdateOrder(&order)

	for _, user := range userInfoList {
		models.AddUser2Room(user.ID, room.ID)
	}
	log.Printf(" [x] Done")

}
