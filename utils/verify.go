package utils

import (
	"fmt"
	"math/rand"

	"github.com/nikhil-thorat/goauth/configs"
	"gopkg.in/gomail.v2"
)

func GenerateVerificationCode() string {

	r := rand.New(rand.NewSource(rand.Int63()))

	code := r.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func SendVerificationEmail(email string, verificationCode string) error {
	mail := gomail.NewMessage()

	mail.SetHeader("From", configs.Envs.EmailFrom)
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "GoAuth Verification Code")

	body := fmt.Sprintf("Your verification code is <b>%v</b>", verificationCode)
	mail.SetBody("text/html", body)

	dialer := gomail.NewDialer(configs.Envs.SmtpHost, int(configs.Envs.SmtpPort), configs.Envs.EmailFrom, configs.Envs.EmailPass)

	err := dialer.DialAndSend(mail)
	if err != nil {
		return fmt.Errorf("FAILED TO SEND VERIFICATION EMAIL : %v", err)
	}

	return nil
}
