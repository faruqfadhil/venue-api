package venue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/faruqfadhil/venue-api/core/entity"
	repoInterface "github.com/faruqfadhil/venue-api/core/repository"
	errutil "github.com/faruqfadhil/venue-api/pkg/error"
	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

const (
	tokenSecretKey = "secret-sekali"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repoInterface.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var out entity.User
	err := r.db.Table("auth").
		Where("LOWER(email) = ?", strings.ToLower(email)).
		First(&out).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("[FindUserByEmail] err: %v", err))
		}
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[FindUserByEmail] err: %v", err))
	}
	return &out, nil
}

func (r *repository) Register(ctx context.Context, payload *entity.User) error {
	err := r.db.Table("auth").Create(&payload).Error
	if err != nil {
		return errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[Register] err: %v", err))
	}
	return nil
}

type jwtClaim struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

func (r *repository) Login(ctx context.Context, email, password string) (*entity.Auth, error) {
	var out entity.User
	err := r.db.Table("auth").
		Where("LOWER(email) = ?", strings.ToLower(email)).
		Where("password = ?", password).
		First(&out).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("[Login] err: %v", err))
		}
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[Login] err: %v", err))
	}

	// Generate access token
	expirationTime := time.Now().Add(24 * time.Hour)
	claim := &jwtClaim{
		out.ID,
		out.Email,
		out.FullName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		}}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	newTokenString, err := newToken.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[Login] err: %v", err))
	}

	return &entity.Auth{
		Email:       out.Email,
		FullName:    out.FullName,
		AccessToken: newTokenString,
	}, nil
}

func (r *repository) ValidateToken(ctx context.Context, token string) (*entity.CredentialClaim, error) {
	claim := &jwtClaim{}
	jwtToken, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecretKey), nil
	})
	if err != nil {
		return nil, errutil.New(errutil.ErrUnauthorized, fmt.Errorf("[ValidateToken] err: %v", err))
	}
	if !jwtToken.Valid {
		return nil, errutil.New(errutil.ErrUnauthorized, fmt.Errorf("[ValidateToken] err: %v", "invalid token"))
	}
	return &entity.CredentialClaim{
		ID:       claim.ID,
		Email:    claim.Email,
		FullName: claim.FullName,
	}, nil
}

func (r *repository) GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error) {

	return nil, nil, nil
}

func (r *repository) GetCities(ctx context.Context) ([]*entity.City, error) {
	var cities []*entity.City
	err := r.db.Table("city").Find(&cities).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetCities] err: %v", err))
	}
	return cities, nil
}
