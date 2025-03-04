package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	// 生成 bcrypt 哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("密码加密失败:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash 校验密码是否与哈希匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
