package rest

import "github.com/gin-gonic/gin"

type Controller interface {
	Map(router *gin.Engine)
}
