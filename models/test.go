package models

import (
	"fmt"
	"go-gin-example/pkg/util"
	"math/rand"
)

func FillUserIfEmpty(nums int) {
	if int(GetUserCount("")) == 0 {
		fmt.Printf("No data in user list, random add %v \n", nums)
		for i := 0; i < nums; i++ {
			tempUser := User{Username: util.RandomString(rand.Int()%16 + 4), Password: util.RandomString(64), Salt: util.RandomString(8)}
			AddUserItem(&tempUser)
		}
	}
}
