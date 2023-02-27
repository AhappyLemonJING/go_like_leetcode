package test

import (
	"crypto/tls"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "wzj <wzj2010624@163.com>"
	e.To = []string{"892263307@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("你的验证码是：<b>123</b>")
	// err := e.Send("smtp.163.com:587", smtp.PlainAuth("", "wzj2010624@163.com", "981122@Wzj", "smtp.163.com"))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// 注意关闭ssl
	err := e.SendWithTLS("smtp.163.com:587", smtp.PlainAuth("", "wzj2010624@163.com", "RUSCZFDRNLMUYJZA", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
