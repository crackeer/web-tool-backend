package api

import (
	"net/http"
	"web-tool-backend/container"

	"github.com/gin-gonic/gin"
)

func GetCrdList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, container.GetCrdList())
}
