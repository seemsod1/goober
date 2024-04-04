package main

import (
	models "github.com/seemsod1/goober/models/app_models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	client, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, string(m.Content))

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}

}
