package utils

import (
	"../configs"
	"log"
	"net/smtp"
)

func SendMail(recipient string, message string) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		configs.MAIL_USERNAME,
		configs.MAIL_PASSWORD,
		configs.MAIL_SERVER,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		configs.MAIL_SERVER + ":" + configs.MAIL_SERVER_PORT,
		auth,
		configs.MAIL_MAIL_ADDRESS,
		[]string{recipient},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}