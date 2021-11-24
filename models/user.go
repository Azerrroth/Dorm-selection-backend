package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model

	Username string     `gorm:"UNIQUE;NOT NULL;UNIQUE_INDEX:username;" json:"username"`
	Password string     `gorm:"NOT NULL" json:"password"`
	IsDel    bool       `gorm:"NOT NULL;default:0" json:"is_del"`
	Salt     string     `gorm:"NOT NULL" json:"salt"`
	LoginAt  *time.Time `json:"login_at"`

	StudentId *string `json:"student_id"`
	Name      *string `json:"name"`
	Gender    uint    `json:"gender"` // 0 未知, 1 男, 2 女
	Mail      *string `json:"Mail"`
	Authority uint    `gorm:"default:0" json:"authority"`
}

func GetUsers(maps interface{}) (users []User) {
	userDB.Where(maps).Find(&users)

	return
}

func GetUserCount(maps interface{}) (count int64) {
	userDB.Model(&User{}).Where(maps).Count(&count)

	return
}

func GetUserInformation(username string) (user User) {
	userDB.Where("username = ?", username).First(&user)
	return
}

type UserInfo struct {
	ID          uint   `json:"id"`
	StudentId   string `json:"student_id"`
	Name        string `json:"name"`
	Gender      uint   `json:"gender"`
	Mail        string `json:"mail"`
	Authority   uint   `gorm:"default:0" json:"authority"`
	CertifyCode string `json:"certify_code"`
}

func GetUserInformationByUsername(username string) (result UserInfo, err error) {
	userTableName := tablePrefix + "users"
	userCertifyTableName := tablePrefix + "user_certifies"

	// sql := fmt.Sprintf(`%s.student_id as student_id, %s.name as name, %s.gender as gender, %s.mail as mail, %s.authority as authority, %s.certify_code as certify_code`, userTableName, userTableName, userTableName, userTableName, userTableName, userCertifyTableName)
	// userDB.Table(userTableName).Select(sql).Where("%s.username = ?", username).First(&result)
	// println(result.Name)
	// println(result.CertifyCode)

	sql := fmt.Sprintf(`SELECT a.*, b.certify_code FROM 
	(SELECT id, student_id, name, gender, mail, authority FROM %s WHERE username = (?)) as a left join %s as b on a.id = b.user_id`, userTableName, userCertifyTableName)
	err = userDB.Raw(sql, username).Scan(&result).Error

	return
}

func GetUserInformationByStudentID(studentID string) (result UserInfo, err error) {
	userTableName := tablePrefix + "users"
	userCertifyTableName := tablePrefix + "user_certifies"

	sql := fmt.Sprintf(`SELECT a.*, b.certify_code FROM 
	(SELECT id, student_id, name, gender, mail, authority FROM %s WHERE student_id = (?)) as a left join %s as b on a.id = b.user_id`, userTableName, userCertifyTableName)
	err = userDB.Raw(sql, studentID).Scan(&result).Error

	return
}

func AddUser(username string, password string, salt string) bool {
	user := User{Username: username, Password: password, Salt: salt}
	result := userDB.Create(&user)
	return result.RowsAffected != 0
}

func AddUserItem(user *User) bool {
	result := userDB.Create(&user)
	return result.RowsAffected != 0
}

func GetUserByID(id uint) (user User) {
	userDB.First(&user, id)
	return
}

func UpdateUser(user *User) bool {
	result := userDB.Save(&user)
	return result.RowsAffected != 0
}

func ExistUserByUsername(username string) bool {
	var count int64
	userDB.Model(&User{}).Where("username = ?", username).Count(&count)
	return int(count) > 0
}

func ExistUserByID(id uint) bool {
	var count int64
	userDB.Model(&User{}).Where("id = ?", id).Count(&count)
	return int(count) > 0
}

func ExistUserByUsernameAndPassword(username string, password string) (User, bool) {
	if !ExistUserByUsername(username) {
		return User{}, false
	}
	var user User
	userDB.Model(&User{}).Where("username = ?", username).First(&user)
	// encrypted, _ := util.EncryptWithSalt(password, []byte(user.Salt))
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt))
	if err != nil {
		return user, false
	} else {
		return user, true
	}
}
