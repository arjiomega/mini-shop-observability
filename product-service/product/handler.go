package product

import (
	"context"
	"errors"
	productpb "mini-shop/proto/productpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	productpb.UnimplementedProductServiceServer

	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateProduct(
	ctx context.Context,
	req *productpb.CreateProductRequest,
) (*productpb.CreateProductResponse, error) {

	product, err := h.service.Create(ProductCreate{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}

	return &productpb.CreateProductResponse{
		Id:   int64(product.ID),
		Name: product.Name,
	}, nil
}
func (h *Handler) GetProduct(
	ctx context.Context,
	req *productpb.GetProductRequest,
) (*productpb.GetProductResponse, error) {

	product, err := h.service.GetByID(int(req.Id))
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"product not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"failed to get product",
		)
	}

	return &productpb.GetProductResponse{
		Id:   int64(product.ID),
		Name: product.Name,
	}, nil
}
