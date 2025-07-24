package email

import (
    "fmt"
    "net/smtp"
)

type Sender interface {
    Send(to, subject, body string) error
}

type SMTPSender struct {
    From     string
    Password string
    Host     string
    Port     string
}

func NewSMTPSender(from, password, host, port string) *SMTPSender {
    return &SMTPSender{From: from, Password: password, Host: host, Port: port}
}

func (s *SMTPSender) Send(to, subject, body string) error {
    addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
    auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "MIME-version: 1.0;\r\n" +
        "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
        body + "\r\n")

    return smtp.SendMail(addr, auth, s.From, []string{to}, msg)
}
