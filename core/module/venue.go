package module

import (
	"context"

	"github.com/faruqfadhil/venue-api/core/entity"
	"github.com/faruqfadhil/venue-api/core/repository"
)

type Usecase interface {
	GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error)
	GetCities(ctx context.Context) ([]*entity.City, error)
}

type usecase struct {
	repo repository.Repository
}

func New(repo repository.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error) {
	return nil, nil, nil
}

func (u *usecase) GetCities(ctx context.Context) ([]*entity.City, error) {
	return nil, nil
}
