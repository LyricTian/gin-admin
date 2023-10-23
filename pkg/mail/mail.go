package mail

import (
	"context"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	globalSender *SmtpSender
	once         sync.Once
)

// Set a global SMTP sender
func SetSender(sender *SmtpSender) {
	once.Do(func() {
		globalSender = sender
	})
}

// Use smtp client send email with to/cc/bcc
func Send(ctx context.Context, to []string, cc []string, bcc []string, subject string, body string, file ...string) error {
	return globalSender.Send(ctx, to, cc, bcc, subject, body, file...)
}

// Use smtp client send email, use to specify recipients
func SendTo(ctx context.Context, to []string, subject string, body string, file ...string) error {
	return globalSender.SendTo(ctx, to, subject, body, file...)
}

// A smtp email client
type SmtpSender struct {
	SmtpHost string
	Port     int
	FromName string
	FromMail string
	UserName string
	AuthCode string
}

func (s *SmtpSender) Send(ctx context.Context, to []string, cc []string, bcc []string, subject string, body string, file ...string) error {
	msg := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	msg.SetHeader("From", msg.FormatAddress(s.FromMail, s.FromName))
	msg.SetHeader("To", to...)
	msg.SetHeader("Cc", cc...)
	msg.SetHeader("Bcc", bcc...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html;charset=utf-8", body)

	for _, v := range file {
		msg.Attach(v)
	}

	d := gomail.NewDialer(s.SmtpHost, s.Port, s.UserName, s.AuthCode)
	return d.DialAndSend(msg)
}

func (s *SmtpSender) SendTo(ctx context.Context, to []string, subject string, body string, file ...string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = s.Send(ctx, to, nil, nil, subject, body, file...)
		if err != nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		err = nil
		break
	}
	return err
}
