package venue

import (
	"context"

	"github.com/faruqfadhil/venue-api/core/entity"
	repoInterface "github.com/faruqfadhil/venue-api/core/repository"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repoInterface.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error) {
	return nil, nil, nil
}
func (r *repository) GetCities(ctx context.Context) ([]*entity.City, error) {
	return nil, nil
}
