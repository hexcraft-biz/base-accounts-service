package config

import (
	"errors"
	"os"

	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
)

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
//
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
