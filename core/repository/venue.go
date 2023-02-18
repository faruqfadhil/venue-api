package repository

import (
	"context"

	"github.com/faruqfadhil/venue-api/core/entity"
)

type Repository interface {
	Register(ctx context.Context, payload *entity.User) error
	Login(ctx context.Context, email, password string) (*entity.Auth, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.CredentialClaim, error)

	GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error)
	GetCities(ctx context.Context) ([]*entity.City, error)
}
