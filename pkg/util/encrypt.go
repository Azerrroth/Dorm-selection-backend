package util

import (
	"errors"
	"go-gin-example/pkg/setting"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Number 随机生成size个数字
func Number(size int) []byte {
	if size <= 0 || size > 10 {
		size = 10
	}
	warehouse := []int{48, 57}
	result := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result[i] = uint8(warehouse[0] + rand.Intn(9))
	}
	return result
}

// Lower 随机生成size个小写字母
func Lower(size int) []byte {
	if size <= 0 || size > 26 {
		size = 26
	}
	warehouse := []int{97, 122}
	result := make([]byte, 26)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result[i] = uint8(warehouse[0] + rand.Intn(26))
	}
	return result
}

// Lower 随机生成size个小写字母
func Upper(size int) []byte {
	if size <= 0 || size > 26 {
		size = 26
	}
	warehouse := []int{65, 90}
	result := make([]byte, 26)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result[i] = uint8(warehouse[0] + rand.Intn(26))
	}
	return result
}

// Salt 生成一个盐值
func Salt(size int) (string, error) {
	// 参数校验
	if size < 0 {
		return "", errors.New("非法的长度")
	}

	// 按需要生成字符串
	var result string
	rand.Seed(time.Now().UnixNano())
	var i int
	for i = 0; i < size; i += 1 {
		randInt := rand.Intn(3)

		switch randInt {
		case 0:
			tempByte := string(Lower(1))
			result += tempByte
		case 1:
			tempByte := string(Number(1))
			result += tempByte
		case 2:
			tempByte := string(Upper(1))
			result += tempByte
		}
	}

	return result, nil
}

func loadConfig() (saltBytes int, hashBytes int) {
	sec, err := setting.Cfg.GetSection("password")
	if err != nil {
		log.Fatal(2, "Fail to get section 'password': %v", err)
	}

	saltBytes, _ = sec.Key("SALT_BYTES").Int()
	hashBytes, _ = sec.Key("HASH_BYTES").Int()
	return
}

func Encrypt(password string) (passwd []byte, salt []byte, errorMsg error) {
	saltBytes, _ := loadConfig()
	saltStr, _ := Salt(saltBytes)
	salt = []byte(saltStr)
	passwd, errorMsg = bcrypt.GenerateFromPassword(append([]byte(password), salt...), bcrypt.DefaultCost)
	return
}

func EncryptWithSalt(password string, salt []byte) (passwd []byte, errorMsg error) {
	// _, hashBytes := loadConfig()
	passwd, errorMsg = bcrypt.GenerateFromPassword(append([]byte(password), salt...), bcrypt.DefaultCost)
	return
}
