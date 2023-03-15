package controllers

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/hexcraft-biz/base-accounts-service/config"
	"github.com/hexcraft-biz/base-accounts-service/misc"
	"github.com/hexcraft-biz/base-accounts-service/models"
	"github.com/hexcraft-biz/controller"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_STATUS_ENABLED            = "enabled"
	EMAIL_CONFIRMATION_EXPIRE_MINS = 10
	JWT_TYPE_SIGN_UP               = "signup"
	JWT_TYPE_FORGET_PWD            = "forgetpwd"
)

type Auth struct {
	*controller.Prototype
	Config config.ConfigInterface
}

func NewAuth(cfg config.ConfigInterface) *Auth {
	return &Auth{
		Prototype: controller.New("auth", cfg.GetDB()),
		Config:    cfg,
	}
}

// ================================================================
// Auth Login
// ================================================================
type genTokenParams struct {
	Identity string `json:"identity" binding:"required,email,min=1,max=128"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params genTokenParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Identity); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				saltedPwd := append([]byte(params.Password), entityRes.Salt...)
				compareErr := bcrypt.CompareHashAndPassword(entityRes.Password, saltedPwd)
				if compareErr != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Password is wrong."})
					return
				}

				if entityRes.Status != USER_STATUS_ENABLED {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "This account is not enabled."})
					return
				}

				if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusOK, absRes)
					return
				}
			}
		}
	}
}

// ================================================================
// SignUp
// ================================================================
type signUpEmailConfirmParams struct {
	Email         string `json:"email" binding:"required,email,min=1,max=128"`
	VerifyPageUrl string `json:"verifyPageURL" binding:"required,url"`
	Continue      string `json:"continue" binding:"omitempty,url"`
}

type signUpEmailConfirmResp struct {
	Token string `json:"token"`
}

func (ctrl *Auth) SignUpEmailConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			params signUpEmailConfirmParams
			uri    *url.URL
		)
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} else if uri, err = url.ParseRequestURI(params.VerifyPageUrl); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
			return
		} else if entityRes != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "This Email is already exist."})
			return
		}

		nowTime := time.Now()
		expiresAt := nowTime.Add(EMAIL_CONFIRMATION_EXPIRE_MINS * time.Minute).Unix()
		issuedAt := nowTime.Unix()

		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		tokenString, err := miscJWT.GenToken(jwt.SigningMethodHS512, misc.EmailJwtClaims{
			StandardClaims: jwt.StandardClaims{
				Subject:   params.Email,
				ExpiresAt: expiresAt,
				IssuedAt:  issuedAt,
			},
			Email:    params.Email,
			Type:     JWT_TYPE_SIGN_UP,
			Continue: params.Continue,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		vals := uri.Query()
		vals.Add("token", tokenString)
		realVerifyPageURI := uri.Scheme + "://" + uri.Host + uri.Path + "?" + vals.Encode()

		tmpl, _ := template.New("email").Parse(getEmailTplHTML())
		var tpl bytes.Buffer

		tmpl.Execute(&tpl, struct {
			Content           string
			RealVerifyPageURI string
			LinkText          string
		}{
			ctrl.Config.GetSignupEmailContent(),
			realVerifyPageURI,
			ctrl.Config.GetSignupEmailLinkText(),
		})

		// TODO Supports multi languages.
		email := misc.NewEmail(
			ctrl.Config.GetSMTPHost(),
			ctrl.Config.GetSMTPPort(),
			ctrl.Config.GetSMTPUsername(),
			ctrl.Config.GetSMTPPassword(),
		)
		email.SendHTML(
			ctrl.Config.GetSMTPSenderName(),
			ctrl.Config.GetSMTPSender(),
			[]string{params.Email},
			ctrl.Config.GetSignupEmailSubject(),
			tpl.String(),
		)

		c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": http.StatusText(http.StatusAccepted)})
		return
	}
}

type signUpTokenVerifyParams struct {
	Token string `form:"token" binding:"required"`
}

type signUpTokenVerifyResp struct {
	Email    string `json:"email"`
	Continue string `json:"continue"`
}

func (ctrl *Auth) SignUpTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signUpTokenVerifyParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_SIGN_UP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, signUpTokenVerifyResp{
			Email:    claims.Email,
			Continue: claims.Continue,
		})
		return
	}
}

type signupParams struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signupParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_SIGN_UP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// TODO Enhanced password requirements.
		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).Insert(claims.Email, params.Password, USER_STATUS_ENABLED); err != nil {
			if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusCreated, absRes)
				return
			}
		}
	}
}

// ================================================================
// ForgetPassword
// ================================================================
type forgetPwdConfirmParams struct {
	Email         string `json:"email" binding:"required,email,min=1,max=128"`
	VerifyPageUrl string `json:"verifyPageURL" binding:"required,url"`
	Continue      string `json:"continue" binding:"omitempty,url"`
}

type forgetPwdConfirmResp struct {
	Token string `json:"token"`
}

func (ctrl *Auth) ForgetPwdConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			params forgetPwdConfirmParams
			uri    *url.URL
		)

		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} else if uri, err = url.ParseRequestURI(params.VerifyPageUrl); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
			return
		} else if entityRes == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "This Email is not already exist."})
			return
		}

		nowTime := time.Now()
		expiresAt := nowTime.Add(EMAIL_CONFIRMATION_EXPIRE_MINS * time.Minute).Unix()
		issuedAt := nowTime.Unix()

		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		tokenString, err := miscJWT.GenToken(jwt.SigningMethodHS512, misc.EmailJwtClaims{
			StandardClaims: jwt.StandardClaims{
				Subject:   params.Email,
				ExpiresAt: expiresAt,
				IssuedAt:  issuedAt,
			},
			Email:    params.Email,
			Type:     JWT_TYPE_FORGET_PWD,
			Continue: params.Continue,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		vals := uri.Query()
		vals.Add("token", tokenString)
		realVerifyPageURI := uri.Scheme + "://" + uri.Host + uri.Path + "?" + vals.Encode()

		tmpl, _ := template.New("email").Parse(getEmailTplHTML())
		var tpl bytes.Buffer

		tmpl.Execute(&tpl, struct {
			Content           string
			RealVerifyPageURI string
			LinkText          string
		}{
			ctrl.Config.GetForgetPwdEmailContent(),
			realVerifyPageURI,
			ctrl.Config.GetForgetPwdEmailLinkText(),
		})

		// TODO Supports multi languages.
		email := misc.NewEmail(
			ctrl.Config.GetSMTPHost(),
			ctrl.Config.GetSMTPPort(),
			ctrl.Config.GetSMTPUsername(),
			ctrl.Config.GetSMTPPassword(),
		)
		email.SendHTML(
			ctrl.Config.GetSMTPSenderName(),
			ctrl.Config.GetSMTPSender(),
			[]string{params.Email},
			ctrl.Config.GetForgetPwdEmailSubject(),
			tpl.String(),
		)

		c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"message": http.StatusText(http.StatusAccepted)})
		return
	}
}

type forgetPwdTokenVerifyParams struct {
	Token string `form:"token" binding:"required"`
}

type forgetPwdTokenVerifyResp struct {
	Email    string `json:"email"`
	Continue string `json:"continue"`
}

func (ctrl *Auth) ForgetPwdTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params forgetPwdTokenVerifyParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_FORGET_PWD {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, forgetPwdTokenVerifyResp{
			Email:    claims.Email,
			Continue: claims.Continue,
		})
		return
	}
}

type forgetPwdParams struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params forgetPwdParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_FORGET_PWD {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		usersEngine := models.NewUsersTableEngine(ctrl.DB)

		if entityRes, err := usersEngine.GetByIdentity(claims.Email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				// TODO next version about password log
				saltedPwd := append([]byte(params.Password), entityRes.Salt...)
				compareErr := bcrypt.CompareHashAndPassword(entityRes.Password, saltedPwd)
				if compareErr == nil {
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
					return
				}

				if _, err := usersEngine.ResetPwd(entityRes.ID, params.Password, entityRes.Salt); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
					return
				}

			}
		}
	}
}

func getEmailTplHTML() string {
	return `
	<!DOCTYPE html>
		<html>
			<body>
				<div>
					<p>{{ .Content }}</p>
					<a href={{ .RealVerifyPageURI }}>{{ .LinkText }}</a>
				</div>
			</body>
		</html>`
}
