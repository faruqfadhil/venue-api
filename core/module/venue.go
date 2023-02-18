package module

import (
	"context"
	"errors"

	"github.com/faruqfadhil/venue-api/core/entity"
	"github.com/faruqfadhil/venue-api/core/repository"
	errutil "github.com/faruqfadhil/venue-api/pkg/error"
)

type Usecase interface {
	GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error)
	GetCities(ctx context.Context) ([]*entity.City, error)
	Register(ctx context.Context, payload *entity.User) error
	Login(ctx context.Context, email, password string) (*entity.Auth, error)
	ValidateToken(ctx context.Context, token string) (*entity.CredentialClaim, error)
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
	return u.repo.GetCities(ctx)
}

func (u *usecase) Register(ctx context.Context, payload *entity.User) error {
	existingUser, err := u.repo.FindUserByEmail(ctx, payload.Email)

	if err != nil && !errors.Is(errutil.GetTypeErr(err), errutil.ErrGeneralNotFound) {
		return err
	}
	if existingUser != nil {
		return errutil.New(errutil.ErrGeneralBadRequest, err, "Email sudah terdaftar di sistem")
	}
	return u.repo.Register(ctx, payload)
}

func (u *usecase) Login(ctx context.Context, email, password string) (*entity.Auth, error) {
	authInfo, err := u.repo.Login(ctx, email, password)
	if err != nil {
		if errors.Is(errutil.GetTypeErr(err), errutil.ErrGeneralNotFound) {
			// Unauthorized.
			return nil, errutil.New(errutil.ErrUnauthorized, err, "Username atau password salah")
		}
		return nil, err
	}
	return authInfo, nil
}

func (u *usecase) ValidateToken(ctx context.Context, token string) (*entity.CredentialClaim, error) {
	return u.repo.ValidateToken(ctx, token)
}
