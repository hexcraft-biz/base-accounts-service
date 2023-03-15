package config

import "github.com/jmoiron/sqlx"

type ConfigInterface interface {
	GetDB() *sqlx.DB
	GetTrustProxy() string
	GetJWTSecret() []byte
	GetSMTPHost() string
	GetSMTPPort() string
	GetSMTPUsername() string
	GetSMTPPassword() string
	GetSMTPSender() string
	GetSMTPSenderName() string
	GetSignupEmailSubject() string
	GetSignupEmailContent() string
	GetSignupEmailLinkText() string
	GetForgetPwdEmailSubject() string
	GetForgetPwdEmailContent() string
	GetForgetPwdEmailLinkText() string
}
