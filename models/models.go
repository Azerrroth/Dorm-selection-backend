package models

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-gin-example/pkg/gender"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
)

// var db *gorm.DB

var userDB, dormDB, orderDB *gorm.DB
var tablePrefix string

var redisDB *redis.Client
var ctx = context.Background()

type Model struct {
	// gorm.Model
	ID        uint      `gorm:"primaryKey;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func init() {
	var (
		errUser, errDorm                           error
		dbName, user, password, userHost, dormHost string
		waitTime, retryTimes                       int
	)
	var redisHost, redisPassword, redisPort string
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	// dbType = sec.Key("TYPE").String() // mysql
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	// host = sec.Key("HOST").String()
	userHost = sec.Key("USER_HOST").String()
	dormHost = sec.Key("DORM_HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()
	waitTime, _ = sec.Key("WAIT_TIME").Int()
	retryTimes, _ = sec.Key("RETRY_TIMES").Int()
	redisHost = sec.Key("REDIS_HOST").String()
	redisPassword = sec.Key("REDIS_PASSWORD").String()
	redisPort = sec.Key("REDIS_PORT").String()

	rand.Seed(time.Now().Unix())
	// db, err = ConnectDB(user, password, host, dbName, tablePrefix)
	userDB, errUser = ConnectDB(user, password, userHost, dbName, tablePrefix)
	dormDB, errDorm = ConnectDB(user, password, dormHost, dbName, tablePrefix)
	if errUser != nil || errDorm != nil {
		fmt.Println(err)
		for i := 0; i < retryTimes; i = i + 1 {
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
			if errUser != nil {
				userDB, errUser = ConnectDB(user, password, userHost, dbName, tablePrefix)
			}
			if errDorm != nil {
				dormDB, errDorm = ConnectDB(user, password, dormHost, dbName, tablePrefix)
			}
			fmt.Printf("Error: connect error, retry times: %d/%d. \n", i, retryTimes)
			if errUser == nil && errDorm == nil {
				break
			}
		}
	}

	fmt.Println(redisHost, redisPort, redisPassword)
	redisDB = ConnectRedis(redisHost, redisPort, redisPassword, retryTimes, time.Duration(waitTime)*time.Millisecond)

	err = redisDB.Set(ctx, "test", "test11111", time.Second*time.Duration(60)).Err()
	if err != nil {
		log.Fatal(2, "Fail to set redis: %v", err)
	}

	val, err := redisDB.Get(ctx, "test").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	orderDB = userDB

	userDB.AutoMigrate(&User{})
	userDB.AutoMigrate(&UserCertify{})

	dormDB.AutoMigrate(&User2Room{})
	dormDB.AutoMigrate(&Building{})
	dormDB.AutoMigrate(&Room{})

	orderDB.AutoMigrate(&Order{})
	orderDB.AutoMigrate(&OrderDetail{})

	// FillDormIfEmpty(fillNums)
	// FillUserIfEmpty(fillNums)

	if initDatabase := os.Getenv("INIT_DB_IF_EMPTY"); initDatabase != "" {
		if int(GetUserCount("")) == 0 {
			initializeData()
		}
	}
	if err != nil {
		log.Println(err)
	}

}

func ConnectDB(user string, password string, host string, dbName string, tablePrefix string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
	})
}

// func CloseDB() {
// 	defer db.Close()
// }

func Test() {
	// AddOrderDetail(&OrderDetail{ResidentId: 1, OrderId: 1})
	// AddOrderDetail(&OrderDetail{ResidentId: 2, OrderId: 1})
	// AddOrderDetail(&OrderDetail{ResidentId: 3, OrderId: 1})

	// result, _ := GetBuildings()

	// for _, userInfo := range result {
	// 	fmt.Println(userInfo.ID)
	// 	fmt.Println(userInfo.Name)
	// }

	// for i := 1; i <= 60; i++ {
	// 	var numstr string
	// 	if i < 10 {
	// 		numstr = fmt.Sprintf("0%d", i)
	// 	} else {
	// 		numstr = fmt.Sprintf("%d", i)
	// 	}
	// 	AddRoom(&Room{BuildingId: 1, Name: fmt.Sprintf("1%s", numstr), Gender: gender.MALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
	// 	AddRoom(&Room{BuildingId: 1, Name: fmt.Sprintf("2%s", numstr), Gender: gender.MALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
	// 	AddRoom(&Room{BuildingId: 1, Name: fmt.Sprintf("3%s", numstr), Gender: gender.FEMALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
	// 	AddRoom(&Room{BuildingId: 1, Name: fmt.Sprintf("4%s", numstr), Gender: gender.FEMALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
	// }

}

func initializeData() {
	buildings := 10
	rooms := 50
	users := 50
	fmt.Printf("Init user: %v \n", users)
	InitUsers(users)

	fmt.Printf("Init buildings: %v \n", buildings)
	InitBuildings(buildings)

	fmt.Printf("Init rooms for every buildings, rooms num per floor: %v \n", rooms)
	InitRooms(buildings, rooms)
}

func InitBuildings(nums int) {
	for i := 1; i <= nums; i++ {
		AddBuilding(fmt.Sprintf("%v号楼", i))
	}
}

func InitRooms(buildingNums int, roomNumsPerFloor int) {
	for buildId := 1; buildId <= buildingNums; buildId++ {
		for i := 1; i <= roomNumsPerFloor; i++ {
			var numstr string
			if i < 10 {
				numstr = fmt.Sprintf("0%d", i)
			} else {
				numstr = fmt.Sprintf("%d", i)
			}
			buildingStr := fmt.Sprintf("%v号楼", buildId)
			AddRoom(&Room{BuildingId: buildId, Name: fmt.Sprintf("%s 1%s", buildingStr, numstr), Gender: gender.MALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
			AddRoom(&Room{BuildingId: buildId, Name: fmt.Sprintf("%s 2%s", buildingStr, numstr), Gender: gender.MALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
			AddRoom(&Room{BuildingId: buildId, Name: fmt.Sprintf("%s 3%s", buildingStr, numstr), Gender: gender.FEMALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
			AddRoom(&Room{BuildingId: buildId, Name: fmt.Sprintf("%s 4%s", buildingStr, numstr), Gender: gender.FEMALE, TotalBeds: 4, AvailableBeds: uint(rand.Intn(5))})
		}
	}
}

func InitUsers(nums int) {
	password := "123456"
	for i := 1; i <= nums; i++ {
		crypted, salt, _ := util.Encrypt(password)
		AddUser(fmt.Sprintf("testuser%d", i), string(crypted), string(salt))
	}
}
