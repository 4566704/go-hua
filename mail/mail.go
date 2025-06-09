package mail

import (
	"net/smtp"
	"strings"
)

func SendToMail(user, password, host, from, subject, body, mailType, replyToAddress string, to, cc, bcc []string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	ccAddress := strings.Join(cc, ";")
	bccAddress := strings.Join(bcc, ";")
	toAddress := strings.Join(to, ";")
	msg := []byte("To: " + toAddress + "\r\nFrom: " + from + "\r\nSubject: " + subject + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + ccAddress + "\r\nBcc: " + bccAddress + "\r\n" + contentType + "\r\n\r\n" + body)

	sendTo := MergeSlice(to, cc)
	sendTo = MergeSlice(sendTo, bcc)
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
