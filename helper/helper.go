package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成Md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(identity string, name string, isAdim int) (string, error) {
	userClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdim,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse token error:%v", err)
	}
	return userClaim, nil
}

// 发送验证码
func SendCode(toUserEmal, code string) error {
	e := email.NewEmail()
	e.From = "wzj <wzj2010624@163.com>"
	e.To = []string{toUserEmal}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("你的验证码是：<b>" + code + "</b>")
	err := e.SendWithTLS("smtp.163.com:587", smtp.PlainAuth("", "wzj2010624@163.com", "RUSCZFDRNLMUYJZA", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	return err
}

// 生成uuid
func GetUUID() string {
	return uuid.NewV4().String()
}

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s = s + strconv.Itoa(rand.Intn(10))
	}
	return s
}

// 代码保存
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}
