package product

import (
	"log"
	"net/http"

	"mini-shop/gateway/api"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client ProductClient
}

func NewHandler(client ProductClient) *Handler {
	return &Handler{
		client: client,
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
		log.Printf("GetProduct gRPC error: %v", err)

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "product service unavailable",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}
