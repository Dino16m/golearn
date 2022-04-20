package routers

import (
	"golearn-api-template/container"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, app container.App) {
	registerAPI(r, app)
	registerWeb(r, app)
}
