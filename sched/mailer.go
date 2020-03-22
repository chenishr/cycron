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
	MailTo  []string
	Subject string
	Body    string
}

var (
	GMailer  *Mailer
	mailChan chan *MailData
)

func init() {
	GMailer = &Mailer{}

	mailChan = make(chan *MailData, 100)

	go func() {
		for {
			select {
			case m, ok := <-mailChan:
				if !ok {
					return
				}

				if err := GMailer.send(m.MailTo, m.Subject, m.Body); err != nil {
					fmt.Println("邮件发送失败:", err.Error())
				}
			}
		}
	}()
}

func (m *Mailer) send(mailTo []string, subject string, body string) error {
	var (
		// 邮件配置
		mailConf conf.MailConf
	)

	// 邮件配置
	mailConf = conf.GConfig.Mail

	mail := gomail.NewMessage()

	//这种方式可以添加别名
	mail.SetHeader("From", "Cycron"+"<"+mailConf.User+">")

	//发送给多个用户
	mail.SetHeader("To", mailTo...)

	//设置邮件主题
	mail.SetHeader("Subject", subject)

	//设置邮件正文
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.User, mailConf.PassWord)

	fmt.Println("发送邮件：" + subject)
	err := d.DialAndSend(mail)
	return err
}

func (m *Mailer) SendMail(mailTo []string, subject, body string) bool {
	var (
		mailData *MailData
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

func (m *Mailer) OrgData(res *ExecResult) {
	var (
		subject string
		body    string
		status  string
		errMsg  string
	)

	psTime := float64(res.endTime.Sub(res.realTime)/time.Millisecond) / 1000
	if res.err == nil {
		status = "【正常】"
	} else {
		status = "【异常】"
	}

	subject = status + "【" + res.job.taskName + "】执行结果"

	if res.err != nil {
		errMsg = `
<p>-------------以下是任务执行错误输出-------------</p>
<p>` + res.err.Error() + `</p>
<p>`
	} else {
		errMsg = ""
	}

	body = `你好，<br/>

<p>以下是任务执行结果：</p>

<p>
	任务ID：` + strconv.FormatInt(res.job.taskId, 10) + `<br/>
	任务名称：` + res.job.taskName + `<br/>
	执行时间：` + res.realTime.Format("2006-01-02 15:04:05") + `<br />
	执行耗时：` + strconv.FormatFloat(psTime, 'g', 6, 64) + `秒<br />
	执行状态：` + status + `
</p>
<p>-------------以下是任务执行输出-------------</p>
<p>` + string(res.output) + `</p>
<p>
` + errMsg + `
--------------------------------------------<br />
	本邮件由系统自动发出，请勿回复<br />
	如果要取消邮件通知，请登录到系统进行设置<br />
</p>`

	m.SendMail(res.job.notifyEmail, subject, body)

	return
}
