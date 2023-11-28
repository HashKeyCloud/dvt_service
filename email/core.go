package email

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"

	"DVT_Service/conf"
)

// Email Mail sending function package
type Email struct {
	dialer *gomail.Dialer

	disable bool

	account string

	to []string
	cc []string
}

func InitEmail(conf *conf.ConfigEmail) *Email {
	e := &Email{
		account: conf.Account,
	}

	if len(conf.To) == 0 {
		e.disable = true
		log.Error().Int("EmailTo", len(e.to)).Msg("disable email, due to empty email receiver")
		return e
	}

	if len(conf.To) == 0 {
		e.disable = true
		log.Error().Msg("disable email, due to To is empty")
		return e
	}
	e.to = conf.To

	if len(conf.CC) > 0 {
		e.cc = conf.CC
	}

	smtp := strings.Split(conf.Smtp, ":")
	if len(smtp) != 2 {
		e.disable = true
		log.Error().Str("smtpHost", conf.Smtp).Msg("disable email, due to error smtpHost")
		return e
	}

	port, err := strconv.Atoi(smtp[1])
	if err != nil {
		e.disable = true
		log.Error().Str("smtpHost", conf.Smtp).Msg("disable email, due to port not number")
		return e
	}

	d := gomail.NewDialer(smtp[0], port, conf.Account, conf.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	e.dialer = d
	log.Info().Msg("Email Open")

	return e
}

// sendToMail send a mail
func (e *Email) sendToMail(subject, content string) error {
	if e.disable {
		return nil
	}

	m := gomail.NewMessage()

	m.SetHeader("From", e.account)

	m.SetHeader("To", e.to...)

	if len(e.cc) > 0 {
		m.SetHeader("Cc", e.cc...)
	}

	m.SetHeader("Subject", subject)

	m.SetBody("text/html; charset=UTF-8", fmt.Sprintf(emailTemplate, content))

	return e.dialer.DialAndSend(m)
}

const emailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Email</title>
</head>
<body>
<div style="height: 5px;width: auto;background-color: red"></div>
<div style="text-align: center;margin-bottom: 20px;font-size: 35px;font-weight: bolder">
    DVT Service
</div>
<div style="display:table;margin:0 auto;">
    %s
</div>
<div style="margin-top: 10px;height: 5px;width: auto;background-color: red"></div>
</body>
</html>
`
