package product

import (
	"log/slog"
	"net/http"

	"mini-shop/gateway/api"
	"mini-shop/gateway/logging"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client ProductClient
	logger *logging.Logger
}

func NewHandler(client ProductClient, logger *logging.Logger) *Handler {
	return &Handler{
		client: client,
		logger: logger,
	}
}

func (h *Handler) GetProduct(
	c *gin.Context,
	productId int,
) {
	product, err := h.client.GetProduct(
		c.Request.Context(),
		int64(productId),
	)
	if err != nil {
		h.logger.Error(
			c.Request.Context(),
			"failed to get product",
			slog.Int("product_id", productId),
			slog.Any("error", err),
		)

		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "product service unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var request api.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn(
			c.Request.Context(),
			"invalid create product request",
			slog.Any("error", err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	product, err := h.client.CreateProduct(
		c.Request.Context(),
		request.Name,
	)
	if err != nil {
		h.logger.Error(
			c.Request.Context(),
			"failed to create product",
			slog.String("product_name", request.Name),
			slog.Any("error", err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "product service unavailable",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}
