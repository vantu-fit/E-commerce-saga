package sender

import (
	"fmt"

	"testing"

	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/saga-pattern/cmd/mail/config"
)

func TestMailBuilder(t *testing.T) {
	builder := NewMailBuilder().
		SetDomain("example.com").
		SetHost("smtp.example.com").
		SetPort(587).
		SetUsername("user@example.com").
		SetPassword("password").
		SetEncryption("tls").
		SetFromAddress("noreply@example.com").
		SetFromName("Example Mailer")

	mail := builder.Build()

	if mail.Domain != "example.com" {
		t.Errorf("Expected Domain to be %s, got %s", "example.com", mail.Domain)
	}

	if mail.Host != "smtp.example.com" {
		t.Errorf("Expected Host to be %s, got %s", "smtp.example.com", mail.Host)
	}

	if mail.Port != 587 {
		t.Errorf("Expected Port to be %d, got %d", 587, mail.Port)
	}

	if mail.Username != "user@example.com" {
		t.Errorf("Expected Username to be %s, got %s", "user@example.com", mail.Username)
	}

	if mail.Password != "password" {
		t.Errorf("Expected Password to be %s, got %s", "password", mail.Password)
	}

	if mail.Encryption != "tls" {
		t.Errorf("Expected Encryption to be %s, got %s", "tls", mail.Encryption)
	}

	if mail.FromAddress != "noreply@example.com" {
		t.Errorf("Expected FromAddress to be %s, got %s", "noreply@example.com", mail.FromAddress)
	}

	if mail.FromName != "Example Mailer" {
		t.Errorf("Expected FromName to be %s, got %s", "Example Mailer", mail.FromName)
	}

	fmt.Println(mail)
}

func TestSendSMTPMessage(t *testing.T) {
	// Tạo một mail builder để tạo đối tượng mail
	// load config
	// load config
	cfgFile, err := config.LoadConfig("../../../cmd/mail/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	fmt.Println(cfg)

	if err != nil {
		log.Error().Msgf("Load config: %v", err)
	}

	// create mail builder
	builder := NewMailBuilder().
		SetDomain("mailhog").
		SetHost("mailhog").
		SetPort(cfg.Mail.MailPortSend).
		SetUsername(cfg.Mail.MailUsername).
		SetPassword(cfg.Mail.MailPassword).
		SetEncryption(cfg.Mail.MailEncryption).
		SetFromAddress(cfg.Mail.MailFromAddress).
		SetFromName(cfg.Mail.MailFromName)

	mail := builder.Build()

	// Tạo một đối tượng Message
	message := Message{
		From:    "norepl@example.com",
		To:      "recipiet@example.com",
		Subject: "Đăng ký thành công",
		Data:    "Do Van Tu",
	}

	// Gửi email và kiểm tra xem có lỗi không
	err = mail.SendRegisterEmail(message)
	require.NoError(t, err)
}
