package interfaces

import (
	"errors"
	"net/http"

	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/shared/resterror"
	"github.com/DimKa163/goseller/internal/user/domain"
	"github.com/DimKa163/goseller/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("seller_email", func(fl validator.FieldLevel) bool {
			email := domain.Email(fl.Field().String())
			return email.Validate() == nil
		})
		_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
			phone := domain.Phone(fl.Field().String())
			return phone.Validate() == nil
		})
	}
}

func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := domain.NewUserID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	user, err := c.app.GetByID(ctx.Request.Context(), userId)
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
	id, err := c.app.Create(ctx.Request.Context(), &req)
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
	us, err := c.app.Update(ctx.Request.Context(), userID, &req)
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
	err = c.app.Delete(ctx.Request.Context(), userId)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *UserController) handleError(ctx *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make([]*resterror.ErrorDetail, len(validationErrors))
		for i, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "required":
				details[i] = &resterror.ErrorDetail{
					Field:   fieldErr.Field(),
					Message: "required",
				}
			case "seller_email":
				details[i] = &resterror.ErrorDetail{
					Field:   fieldErr.Field(),
					Message: "invalid format",
					Value:   fieldErr.Value(),
				}
			case "phone":
				details[i] = &resterror.ErrorDetail{
					Field:   fieldErr.Field(),
					Message: "invalid format",
					Value:   fieldErr.Value(),
				}
			default:
				continue
			}
		}
		ctx.JSON(http.StatusBadRequest, &resterror.ErrorResponse{Error: &resterror.Error{
			Message: "invalid request",
			Details: details,
			Code:    int(shared.ErrorCodeBadInputData),
		}})
		return
	}
	var appErr shared.SellerError
	if errors.As(err, &appErr) {
		details := make([]*resterror.ErrorDetail, len(appErr.Details()))
		for i, d := range appErr.Details() {
			details[i] = &resterror.ErrorDetail{
				Message: d.Message,
			}
		}
		ctx.JSON(shared.ErrorCodeToHttpStatusCode(appErr.GetCode()), &resterror.ErrorResponse{Error: &resterror.Error{
			Message: appErr.Error(),
			Code:    int(appErr.GetCode()),
			Details: details,
		}})
		return
	}
	if errors.Is(err, domain.ErrInvalidUserID) {
		ctx.JSON(http.StatusBadRequest, &resterror.ErrorResponse{Error: &resterror.Error{
			Message: err.Error(),
			Details: []*resterror.ErrorDetail{
				{
					Field:   "id",
					Message: "Invalid user ID format",
				},
			},
			Code: int(shared.ErrorCodeInvalidID),
		}})
		return
	}
	ctx.JSON(http.StatusInternalServerError, &resterror.ErrorResponse{
		Error: &resterror.Error{
			Message: "internal server error",
			Code:    int(shared.ErrorCodeInternalServerError),
		},
	})
}
