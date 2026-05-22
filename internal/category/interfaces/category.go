package interfaces

import (
	"errors"
	"net/http"

	"github.com/DimKa163/goseller/internal/category/domain"
	"github.com/DimKa163/goseller/internal/category/usecase"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/shared/resterror"
	"github.com/DimKa163/goseller/internal/shared/sellerlog"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type CategoryController struct {
	logger *zap.Logger
	app    *usecase.CategoryAppService
}

func NewCategoryController(logger *zap.Logger, app *usecase.CategoryAppService) *CategoryController {
	return &CategoryController{
		logger: logger,
		app:    app,
	}
}

func (c *CategoryController) Map(engine *gin.Engine) {
	engine.GET("category/:id", c.Get)
	engine.POST("category", c.Create)
	engine.PUT("category/:id", c.Update)
	engine.DELETE("category/:id", c.Delete)
}

func (c *CategoryController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryId, err := domain.NewCategoryID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	category, err := c.app.GetByID(ctx.Request.Context(), categoryId)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, category)

}

func (c *CategoryController) Create(ctx *gin.Context) {
	var request usecase.CategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.handleError(ctx, err)
		return
	}
	id, err := c.app.Create(ctx.Request.Context(), &request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Header("Location", "/category/"+id.String())
	ctx.Status(http.StatusCreated)
}

func (c *CategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryId, err := domain.NewCategoryID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	var request usecase.CategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.handleError(ctx, err)
		return
	}
	_, err = c.app.Update(ctx.Request.Context(), categoryId, &request)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *CategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryId, err := domain.NewCategoryID(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	err = c.app.Delete(ctx.Request.Context(), categoryId)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *CategoryController) handleError(ctx *gin.Context, err error) {
	log := sellerlog.FromContext(ctx, c.logger)
	log.Error("error occurred during handling request", zap.Error(err))
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
	var categoryIDErr *domain.CategoryIDInvalidError
	if errors.As(err, &categoryIDErr) {
		ctx.JSON(shared.ErrorCodeToHttpStatusCode(categoryIDErr.GetCode()), &resterror.ErrorResponse{Error: &resterror.Error{
			Message: err.Error(),
			Details: []*resterror.ErrorDetail{
				{
					Field: "id",

					Message: "Invalid category ID format",
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
