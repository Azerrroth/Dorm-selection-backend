package models

import (
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const usernameSet = "usernameSet"
const buildingPrefix = "building_id_"
const expireTime = time.Minute * 60

func ConnectRedis(addr string, port string, password string, retryTimes int, retryBackoff time.Duration) (rds *redis.Client) {
	rds = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", addr, port),
		Password:        password, // no password set
		DB:              0,        // use default DB
		MaxRetries:      retryTimes,
		MinRetryBackoff: retryBackoff,
	})
	return
}

func CheckUsernameIsMember(username string) (isMember bool) {
	isMember = false
	isMember, _ = redisDB.SIsMember(ctx, usernameSet, username).Result()
	return
}

func AddUsernameToSet(username string) (err error) {
	err = redisDB.SAdd(ctx, usernameSet, username).Err()
	return
}

func ClearSet(name string) (err error) {
	total := redisDB.SCard(ctx, name).Val()
	err = redisDB.SPopN(ctx, name, total).Err()
	return
}

func SetBuildingStatus(status BuildingStatus) (err error) {
	str, _ := json.Marshal(status)
	err = redisDB.Set(ctx, fmt.Sprintf(buildingPrefix+"%d", status.BuildingId), str, expireTime).Err()
	return
}

func GetBuildingStatus(buildingId uint) (status BuildingStatus, err error) {
	status = BuildingStatus{}
	var resultStr string
	resultStr, err = redisDB.Get(ctx, fmt.Sprintf(buildingPrefix+"%d", buildingId)).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resultStr), &status)
	return
}

func MinusBuildingStatus(buildingId uint, num int, isMale bool) (err error) {
	status, err := GetBuildingStatus(buildingId)
	if err != nil {
		return
	}
	if isMale {
		status.MaleAvailable += num
	} else {
		status.FemaleAvailable += num
	}
	err = SetBuildingStatus(status)
	return
}
