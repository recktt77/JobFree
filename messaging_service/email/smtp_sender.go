package email

import (
	"net/smtp"
)

type SmtpSender struct {
	from     string
	password string
	host     string
	port     string
}

func NewSmtpSender(from, password, host, port string) *SmtpSender {
	return &SmtpSender{from, password, host, port}
}

func (s *SmtpSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	addr := s.host + ":" + s.port
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
