package interfaces

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryController struct {
	logger *zap.Logger
}

func NewCatetegoryController(logger *zap.Logger) *CategoryController {
	return &CategoryController{
		logger: logger,
	}
}

func (c *CategoryController) Map(engine *gin.Engine) {
	engine.GET("category/:id", c.Get)
	engine.POST("category", c.Create)
	engine.PUT("category/:id", c.Update)
	engine.DELETE("category/:id", c.Delete)
}

func (c *CategoryController) Get(ctx *gin.Context) {
	panic("implement me")
}

func (c *CategoryController) Create(ctx *gin.Context) {
	panic("implement me")
}

func (c *CategoryController) Update(ctx *gin.Context) {
	panic("implement me")
}

func (c *CategoryController) Delete(ctx *gin.Context) {
	panic("implement me")
}
