package repository

import (
	"context"
	"time"

	"github.com/faruqfadhil/venue-api/core/entity"
)

type Repository interface {
	Register(ctx context.Context, payload *entity.User) error
	Login(ctx context.Context, email, password string) (*entity.Auth, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.CredentialClaim, error)

	GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error)
	GetCities(ctx context.Context) ([]*entity.City, error)

	GetOrderByPackageIDAndDate(ctx context.Context, packageID int, date time.Time) (*entity.Order, error)
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetPackageByID(ctx context.Context, ID int) (*entity.VenuePackage, error)
	GetGalleriesByVenueIDs(ctx context.Context, IDs []int) (map[int][]string, error)
	GetOrdersByDate(ctx context.Context, date time.Time) ([]*entity.Order, error)
	GetVenuePackageByQuery(ctx context.Context, param *entity.GetVenuePackageQuery) ([]*entity.VenuePackage, error)
	GetVenueCategoryPackageByQuery(ctx context.Context, param *entity.GetVenueCategoryByQuery) ([]*entity.VenuePackageCategory, error)
}
