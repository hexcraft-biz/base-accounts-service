package main

import (
	"errors"
	"os"

	"github.com/hexcraft-biz/base-accounts-service/service"
	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := Load()
	MustNot(err)
	cfg.DBOpen(false)

	service.New(cfg).Run(":" + cfg.Env.AppPort)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// ================================================================
// Env
// ================================================================
type Env struct {
	*env.Prototype
	JWTSecret              []byte
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	SMTPSender             string
	SignupEmailSubject     string
	SignupEmailContent     string
	SignupEmailLinkText    string
	ForgetPwdEmailSubject  string
	ForgetPwdEmailContent  string
	ForgetPwdEmailLinkText string
}

func FetchEnv() (*Env, error) {
	if e, err := env.Fetch(); err != nil {
		return nil, err
	} else {

		env := &Env{
			Prototype: e,
		}

		if os.Getenv("JWT_SECRET") != "" {
			env.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
		} else {
			return nil, errors.New("Invalid environment variable : JWT_SECRET")
		}

		if os.Getenv("SMTP_HOST") != "" {
			env.SMTPHost = os.Getenv("SMTP_HOST")
		} else {
			return nil, errors.New("Invalid environment variable : SMTP_HOST")
		}

		if os.Getenv("SMTP_PORT") != "" {
			env.SMTPPort = os.Getenv("SMTP_PORT")
		} else {
			return nil, errors.New("Invalid environment variable : SMTP_PORT")
		}

		if os.Getenv("SMTP_USERNAME") != "" {
			env.SMTPUsername = os.Getenv("SMTP_USERNAME")
		} else {
			return nil, errors.New("Invalid environment variable : SMTP_USERNAME")
		}

		if os.Getenv("SMTP_PASSWORD") != "" {
			env.SMTPPassword = os.Getenv("SMTP_PASSWORD")
		} else {
			return nil, errors.New("Invalid environment variable : SMTP_PASSWORD")
		}

		if os.Getenv("SMTP_SENDER") != "" {
			env.SMTPSender = os.Getenv("SMTP_SENDER")
		} else {
			return nil, errors.New("Invalid environment variable : SMTP_SENDER")
		}

		if os.Getenv("SIGNUP_EMAIL_SUBJECT") != "" {
			env.SignupEmailSubject = os.Getenv("SIGNUP_EMAIL_SUBJECT")
		} else {
			return nil, errors.New("Invalid environment variable : SIGNUP_EMAIL_SUBJECT")
		}

		if os.Getenv("SIGNUP_EMAIL_CONTENT") != "" {
			env.SignupEmailContent = os.Getenv("SIGNUP_EMAIL_CONTENT")
		} else {
			return nil, errors.New("Invalid environment variable : SIGNUP_EMAIL_CONTENT")
		}

		if os.Getenv("SIGNUP_EMAIL_LINK_TEXT") != "" {
			env.SignupEmailLinkText = os.Getenv("SIGNUP_EMAIL_LINK_TEXT")
		} else {
			return nil, errors.New("Invalid environment variable : SIGNUP_EMAIL_LINK_TEXT")
		}

		if os.Getenv("FORGET_PWD_EMAIL_SUBJECT") != "" {
			env.ForgetPwdEmailSubject = os.Getenv("FORGET_PWD_EMAIL_SUBJECT")
		} else {
			return nil, errors.New("Invalid environment variable : FORGET_PWD_EMAIL_SUBJECT")
		}

		if os.Getenv("FORGET_PWD_EMAIL_CONTENT") != "" {
			env.ForgetPwdEmailContent = os.Getenv("FORGET_PWD_EMAIL_CONTENT")
		} else {
			return nil, errors.New("Invalid environment variable : FORGET_PWD_EMAIL_CONTENT")
		}

		if os.Getenv("FORGET_PWD_LINK_TEXT") != "" {
			env.ForgetPwdEmailLinkText = os.Getenv("FORGET_PWD_LINK_TEXT")
		} else {
			return nil, errors.New("Invalid environment variable : FORGET_PWD_LINK_TEXT")
		}

		return env, nil
	}
}

// ================================================================
// Config
// ================================================================
type Config struct {
	*Env
	DB *sqlx.DB
}

func Load() (*Config, error) {
	e, err := FetchEnv()
	if err != nil {
		return nil, err
	}

	return &Config{Env: e}, nil
}

func (cfg *Config) DBOpen(init bool) error {
	var err error

	cfg.DBClose()
	cfg.DB, err = cfg.MysqlConnectWithMode(init)

	return err
}

func (cfg *Config) DBClose() {
	if cfg.DB != nil {
		cfg.DB.Close()
	}
}

// ================================================================
// Config implement ConfigInterface
// ================================================================
func (cfg *Config) GetDB() *sqlx.DB {
	return cfg.DB
}

func (cfg *Config) GetTrustProxy() string {
	return cfg.Env.TrustProxy
}

func (cfg *Config) GetJWTSecret() []byte {
	return cfg.Env.JWTSecret
}

func (cfg *Config) GetSMTPHost() string {
	return cfg.Env.SMTPHost
}

func (cfg *Config) GetSMTPPort() string {
	return cfg.Env.SMTPPort
}

func (cfg *Config) GetSMTPUsername() string {
	return cfg.Env.SMTPUsername
}

func (cfg *Config) GetSMTPPassword() string {
	return cfg.Env.SMTPPassword
}

func (cfg *Config) GetSMTPSender() string {
	return cfg.Env.SMTPSender
}

func (cfg *Config) GetSignupEmailSubject() string {
	return cfg.Env.SignupEmailSubject
}

func (cfg *Config) GetSignupEmailContent() string {
	return cfg.Env.SignupEmailContent
}

func (cfg *Config) GetSignupEmailLinkText() string {
	return cfg.Env.SignupEmailLinkText
}

func (cfg *Config) GetForgetPwdEmailSubject() string {
	return cfg.Env.ForgetPwdEmailSubject
}

func (cfg *Config) GetForgetPwdEmailContent() string {
	return cfg.Env.ForgetPwdEmailContent
}

func (cfg *Config) GetForgetPwdEmailLinkText() string {
	return cfg.Env.ForgetPwdEmailLinkText
}
