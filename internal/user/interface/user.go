package interfaces

import (
	"errors"
	"net/http"

	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/transport"
	"github.com/DimKa163/goseller/internal/user/domain"
	"github.com/DimKa163/goseller/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	logger *zap.Logger
	app    *usecase.User
}

func NewUserController(logger *zap.Logger, app *usecase.User) *UserController {
	return &UserController{
		logger: logger,
		app:    app,
	}
}

func (c *UserController) Map(router *gin.Engine) {
	router.GET("/user/:id", c.GetUser)
	router.POST("/user", c.CreateUser)
	router.PUT("/user/:id", c.UpdateUser)
	router.DELETE("/user/:id", c.DeleteUser)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := domain.NewUserID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	user, err := c.app.GetByID(ctx, userId)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req usecase.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.handleError(ctx, err)
		return
	}
	id, err := c.app.Create(ctx, &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Header("Location", "/user/"+id.String())
	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req usecase.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.handleError(ctx, err)
		return
	}
	id := ctx.Param("id")
	userID, err := domain.NewUserID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	us, err := c.app.Update(ctx, userID, &req)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, us)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := domain.NewUserID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	err = c.app.Delete(ctx, userId)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *UserController) handleError(ctx *gin.Context, err error) {
	if errors.Is(err, domain.ErrInvalidUserID) {
		ctx.JSON(http.StatusBadRequest, &transport.ErrorResponse{Error: &transport.Error{
			Message: err.Error(),
			Details: []*transport.ErrorDetail{
				{
					Field:   "id",
					Message: "Invalid user ID format",
				},
			},
			Code: int(shared.ErrorCodeInvalidID),
		}})
		return
	}
	if errors.Is(err, domain.ErrInvalidEmail) {
		ctx.JSON(http.StatusBadRequest, &transport.ErrorResponse{Error: &transport.Error{
			Message: err.Error(),
			Details: []*transport.ErrorDetail{
				{
					Field:   "email",
					Message: "Invalid email format",
				},
			},
			Code: int(shared.ErrorCodeInvalidID),
		}})
		return
	}
	if errors.Is(err, domain.ErrInvalidPhone) {
		ctx.JSON(http.StatusBadRequest, &transport.ErrorResponse{Error: &transport.Error{
			Message: err.Error(),
			Details: []*transport.ErrorDetail{
				{
					Field:   "phone",
					Message: "Invalid phone format",
				},
			},
			Code: int(shared.ErrorCodeInvalidID),
		}})
		return
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, &transport.ErrorResponse{Error: &transport.Error{
			Message: err.Error(),
			Details: []*transport.ErrorDetail{},
			Code:    int(shared.ErrorCodeResourceNotFound),
		}})
		return
	}
	if errors.Is(err, dberror.ErrDuplicateKey) {
		ctx.JSON(http.StatusConflict, &transport.ErrorResponse{Error: &transport.Error{
			Message: err.Error(),
			Details: []*transport.ErrorDetail{},
			Code:    int(shared.ErrorCodeResourceAlreadyExists),
		}})
		return
	}
	ctx.JSON(400, &transport.ErrorResponse{
		Error: &transport.Error{
			Message: "Invalid request body",
			Details: []*transport.ErrorDetail{}},
	})
}
