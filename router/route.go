package router

import (
	"bluebell_renjiexuan/controller"
	"bluebell_renjiexuan/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //将gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/ping ", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1")

	//注册
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	return r
}
