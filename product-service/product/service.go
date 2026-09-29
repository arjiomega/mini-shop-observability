package product

import (
	"context"
	"log"
	"time"
)

type Service struct {
	repository *Repository
	cache      ProductCache
}

func NewService(
	repository *Repository,
	cache ProductCache,
) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
	}
}
func (s *Service) Create(ctx context.Context, input ProductCreate) (Product, error) {
	return s.repository.Create(ctx, input)
}

func (s *Service) GetByID(ctx context.Context, id int) (Product, error) {
	// GET CACHE IF EXISTS
	product, err := s.cache.Get(ctx, id)
	if err == nil {
		return product, nil
	}

	// CACHE MISS
	product, err = s.repository.GetByID(ctx, id)
	if err != nil {
		return Product{}, err
	}

	// POPULATE CACHE
	if err := s.cache.Set(ctx, product, 5*time.Minute); err != nil {
		if ctx.Err() != nil {
			return product, ctx.Err()
		}
		log.Printf("failed to cache product %d: %v", id, err)
	}

	return product, nil
}

func (s *Service) List() ([]Product, error) {
	return s.repository.List()
}
