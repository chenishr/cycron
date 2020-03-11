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
					fmt.Println("邮件发送失败:", err.Error())
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
	mail.SetHeader("From","Cycron" + "<" + mailConn["user"] + ">")

	//发送给多个用户
	mail.SetHeader("To", mailTo...)
	fmt.Println(mailTo)

	//设置邮件主题
	mail.SetHeader("Subject", subject)

	//设置邮件正文
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	fmt.Println("发送邮件：" + subject)
	err := d.DialAndSend(mail)
	return err
}

func (m *Mailer)SendMail(mailTo []string,subject, body string) bool {
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

func (m *Mailer)OrgData(res *ExecResult) {
	var (
		subject	string
		body	string
		status	string
	)

	psTime 	:= float64(res.endTime.Sub(res.info.RealTime) / time.Millisecond) / 1000
	if res.err == nil {
		status = "【正常】"
	}else{
		status = "【异常】"
	}

	subject = status + "【" + res.info.job.taskName + "】执行结果"

	body = `你好，<br/>

<p>以下是任务执行结果：</p>

<p>
	任务ID：` + strconv.FormatInt(int64(res.info.job.taskId),10) + `<br/>
	任务名称：` + res.info.job.taskName + `<br/>
	执行时间：` + res.info.RealTime.Format("2006-01-02 15:04:05")  + `<br />
	执行耗时：` + strconv.FormatFloat(psTime,'g',6,64) + `秒<br />
	执行状态：` + status + `
</p>
<p>-------------以下是任务执行输出-------------</p>
<p>` + string(res.output) + `</p>
<p>
--------------------------------------------<br />
	本邮件由系统自动发出，请勿回复<br />
	如果要取消邮件通知，请登录到系统进行设置<br />
</p>`

	m.SendMail(res.info.job.notifyEmail,subject,body)

	return
}
