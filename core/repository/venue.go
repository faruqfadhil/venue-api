package repository

import (
	"context"

	"github.com/faruqfadhil/venue-api/core/entity"
)

type Repository interface {
	GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error)
	GetCities(ctx context.Context) ([]*entity.City, error)
}
