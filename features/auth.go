package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-accounts-service/config"
	"github.com/hexcraft-biz/base-accounts-service/controllers"
	"github.com/hexcraft-biz/feature"
)

const (
	SCOPE_USER_PROTOTYPE = "user.prototype"
)

func LoadAuth(e *gin.Engine, cfg config.ConfigInterface) {
	c := controllers.NewAuth(cfg)

	authV1 := feature.New(e, "/auth/v1")

	authV1.POST("/signup/confirmation", c.SignUpEmailConfirm())
	authV1.GET("/signup/tokeninfo", c.SignUpTokenVerify())
	authV1.POST("/signup", c.SignUp())

	authV1.POST("/forgetpassword/confirmation", c.ForgetPwdConfirm())
	authV1.GET("/forgetpassword/tokeninfo", c.ForgetPwdTokenVerify())
	authV1.PUT("/password", c.ChangePassword())

	authV1.POST("/login", c.Login())
}
