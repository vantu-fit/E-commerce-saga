package sender


import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// main interface for this features of app
type MailSender interface {
	SendRegisterEmail(message Message) error
}

// main struct of mail sender 
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

type MailBuilder interface {
	SetDomain(domain string) MailBuilder
	SetHost(host string) MailBuilder
	SetPort(port int) MailBuilder
	SetUsername(username string) MailBuilder
	SetPassword(password string) MailBuilder
	SetEncryption(encryption string) MailBuilder
	SetFromAddress(fromAddress string) MailBuilder
	SetFromName(fromName string) MailBuilder
	Build() *Mail
}

// builder pattern
// mailBuilder là triển khai cụ thể của MailBuilder.
type mailBuilder struct {
	domain      string
	host        string
	port        int
	username    string
	password    string
	encryption  string
	fromAddress string
	fromName    string
}

// NewMailBuilder trả về một thể hiện của mailBuilder.
func NewMailBuilder() MailBuilder {
	return &mailBuilder{}
}

// SetDomain thiết lập domain.
func (b *mailBuilder) SetDomain(domain string) MailBuilder {
	b.domain = domain
	return b
}

// SetHost thiết lập host.
func (b *mailBuilder) SetHost(host string) MailBuilder {
	b.host = host
	return b
}

// SetPort thiết lập port.
func (b *mailBuilder) SetPort(port int) MailBuilder {
	b.port = port
	return b
}

// SetUsername thiết lập username.
func (b *mailBuilder) SetUsername(username string) MailBuilder {
	b.username = username
	return b
}

// SetPassword thiết lập password.
func (b *mailBuilder) SetPassword(password string) MailBuilder {
	b.password = password
	return b
}

// SetEncryption thiết lập phương thức mã hóa.
func (b *mailBuilder) SetEncryption(encryption string) MailBuilder {
	b.encryption = encryption
	return b
}

// SetFromAddress thiết lập địa chỉ gửi.
func (b *mailBuilder) SetFromAddress(fromAddress string) MailBuilder {
	b.fromAddress = fromAddress
	return b
}

// SetFromName thiết lập tên người gửi.
func (b *mailBuilder) SetFromName(fromName string) MailBuilder {
	b.fromName = fromName
	return b
}

func (b *mailBuilder) Build() *Mail {
	return &Mail{
		Domain:      b.domain,
		Host:        b.host,
		Port:        b.port,
		Username:    b.username,
		Password:    b.password,
		Encryption:  b.encryption,
		FromAddress: b.fromAddress,
		FromName:    b.fromName,
	}
}

// implement Mail feature

func (m *Mail) SendRegisterEmail(message Message) error {
	if message.From == "" {
		message.From = m.FromAddress
	}
	if message.FromName == "" {
		message.FromName = m.FromName
	}

	data := map[string]any{
		"Name": message.Data,
		"ActivationLink" : "http://localhost:8080/activate",
	}

	message.DataMap = data

	formatMessage, err := m.BuildHTMLMessage(message)
	if err != nil {
		return err
	}

	plainMessage, err := m.BuildPlainTextMessage(message)

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Encryption = m.GetEncryption(m.Encryption)
	server.Username = m.Username
	server.Password = m.Password
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.
		SetFrom(message.From).
		AddTo(message.To).
		SetSubject(message.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formatMessage)

	if len(message.Attachments) > 0 {
		for _, x := range message.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil

}

// build string html
func (m *Mail) BuildHTMLMessage(message Message) (string, error) {
	// for test  "../../templates/mail.html.gohtml"
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", message.DataMap); err != nil {
		return "", err
	}

	formatedMessage := tpl.String()
	formatedMessage, err = m.inlineCss(formatedMessage)
	if err != nil {
		return "", err
	}

	return formatedMessage, nil
}

// build string html no use premailer
func (m *Mail) BuildPlainTextMessage(message Message) (string, error) {
	// for test  "../../templates/mail.html.gohtml"
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", message.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// convert html doc contain css into static html doc
func (m *Mail) inlineCss(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	pre, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := pre.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// string => Encryption
func (m *Mail) GetEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
