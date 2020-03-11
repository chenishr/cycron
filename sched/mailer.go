package sched

import (
	"cycron/conf"
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
	"time"
)

type Mailer struct {

}

type MailData struct {
	MailTo		[]string
	Subject		string
	Body		string
}

var(
	GMailer 	*Mailer
	mailChan	chan *MailData
)

func init()  {
	GMailer = &Mailer{}

	mailChan = make(chan *MailData,100)

	go func() {
		for {
			select {
			case m, ok := <-mailChan:
				if !ok {
					return
				}

				if err := GMailer.send(m.MailTo,m.Subject,m.Body); err != nil {
					fmt.Println("SendMail:", err.Error())
				}
			}
		}
	}()
}

func (m *Mailer)send(mailTo []string,subject string, body string ) error {
	var(
		// 邮件配置
		mailConf 	conf.MailConf
	)

	mailConf = conf.GConfig.Mail

	// 发送邮件配置
	mailConn := map[string]string {
		"user": mailConf.User,
		"pass": mailConf.PassWord,
		"host": mailConf.Host,
		"port": mailConf.Port,
	}

	//转换端口类型为int
	port, _ := strconv.Atoi(mailConn["port"])

	mail := gomail.NewMessage()

	//这种方式可以添加别名
	mail.SetHeader("From","定时任务管理器" + "<" + mailConn["user"] + ">")

	//发送给多个用户
	mail.SetHeader("To", mailTo...)

	//设置邮件主题
	mail.SetHeader("Subject", subject)

	//设置邮件正文
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(mail)
	return err
}

func (m Mailer)SendMail(mailTo []string,subject, body string) bool {
	var (
		mailData	*MailData
	)

	mailData = &MailData{
		MailTo:  mailTo,
		Subject: subject,
		Body:    body,
	}

	select {
	case mailChan <- mailData:
		return true
	case <-time.After(time.Second * 5):
		return false
	}
}
