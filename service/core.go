package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-accounts-service/config"
	"github.com/hexcraft-biz/base-accounts-service/features"
)

func New(cfg config.ConfigInterface) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.GetTrustProxy()})

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg)

	return engine
}
