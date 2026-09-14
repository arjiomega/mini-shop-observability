package product

import (
	"context"
	productpb "mini-shop/proto/productpb"
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
		return nil, err
	}

	return &productpb.GetProductResponse{
		Id:   int64(product.ID),
		Name: product.Name,
	}, nil
}
