package models

import (
	"go-gin-example/pkg/util"
)

const CertifyCodeLength = 6

type UserCertify struct {
	Model

	UserId      uint   `gorm:"" json:"id"`
	CertifyCode string `gorm:"" json:"certify_code"`
}

func GetUserCertify(userId uint) (UserCertify, error) {
	var certify UserCertify
	err := userDB.Where("user_id = ?", userId).First(&certify).Error
	return certify, err
}

func GetCertifyCode(userId uint) (string, error) {
	var certify UserCertify
	err := userDB.Where("user_id = ?", userId).First(&certify).Error
	if err != nil {
		return "", err
	}
	return certify.CertifyCode, nil
}

// Add item in user_certify table, with random certify code.
// Input: username
// Return: certify code, error
func GenerateCertifyCode(username string) (string, error) {
	var user User
	err := userDB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	certifyCode := util.RandomNumber(CertifyCodeLength)
	certify := UserCertify{
		UserId:      user.ID,
		CertifyCode: certifyCode,
	}
	var count int64
	userDB.Model(UserCertify{}).Where("user_id = ?", user.ID).Count(&count)
	if count == 0 {
		err = userDB.Save(&certify).Error
		if err != nil {
			return "", err
		}
		return certifyCode, nil
	}
	return "", err
}

func UpdateCertifyCode(user_id uint) (string, error) {
	var user User
	err := userDB.Where("id = ?", user_id).First(&user).Error
	if err != nil {
		return "", err
	}
	certify_code := util.RandomNumber(CertifyCodeLength)
	certify, err := GetUserCertify(user.ID)
	if err == nil {
		certify.CertifyCode = certify_code
		err = userDB.Save(&certify).Error
		if err != nil {
			return "", err
		}
	} else {
		certify = UserCertify{
			UserId:      user.ID,
			CertifyCode: certify_code,
		}
		err = userDB.Save(&certify).Error
		if err != nil {
			return "", err
		}
	}
	return certify_code, nil
}
