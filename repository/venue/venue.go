package venue

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	var result []Venue
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	qb := r.db.Table("venue")
	if param.ID > 0 {
		qb = qb.Where("id = ?", param.ID)
	}
	if len(param.NotInIDs) > 0 {
		qb = qb.Where("id NOT IN (?)", param.NotInIDs)
	}
	if param.CityID > 0 {
		qb = qb.Where("city_id = ?", param.CityID)
	}
	if len(param.CityIDs) > 0 {
		qb = qb.Where("city_id IN(?)", param.CityIDs)
	}
	if param.IsFavourite {
		qb = qb.Where("is_favourite = ?", param.IsFavourite)
	}

	var pag *entity.Pagination
	if param.IsWithoutPagination {
		err := qb.Find(&result).Error
		if err != nil {
			return nil, nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetVenues] err: %v", err))
		}
	} else {
		var totalRecords int64
		err := qb.Count(&totalRecords).Error
		if err != nil {
			return nil, nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetVenues] err: %v", err))
		}

		offset := (param.Page - 1) * param.Limit
		qb = qb.Order("id asc")
		data := qb
		err = data.Limit(param.Limit).Offset(offset).Find(&result).Error
		if err != nil {
			return nil, nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetVenues] err: %v", err))
		}

		totalPage := math.Ceil(float64(totalRecords) / float64(param.Limit))
		pag = &entity.Pagination{
			Page:         param.Page,
			TotalPage:    int(totalPage),
			CurrentItems: len(result),
			TotalItems:   int(totalRecords),
		}
	}

	out := []*entity.Venue{}
	for _, r := range result {
		out = append(out, r.ToEntity())
	}

	return out, pag, nil
}

func (r *repository) GetCities(ctx context.Context) ([]*entity.City, error) {
	var cities []*entity.City
	err := r.db.Table("city").Find(&cities).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetCities] err: %v", err))
	}
	return cities, nil
}

func (r *repository) GetOrderByPackageIDAndDate(ctx context.Context, packageID int, date time.Time) (*entity.Order, error) {
	var out entity.Order
	err := r.db.Table("order").
		Where("package_id = ?", packageID).
		Where("date = ?", date).
		First(&out).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("[GetOrderByPackageIDAndDate] err: %v", err))
		}
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetOrderByPackageIDAndDate] err: %v", err))
	}
	return &out, nil
}

func (r *repository) CreateOrder(ctx context.Context, order *entity.Order) error {
	err := r.db.Table("order").Create(&order).Error
	if err != nil {
		return errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[CreateOrder] err: %v", err))
	}
	return nil
}

func (r *repository) GetPackageByID(ctx context.Context, ID int) (*entity.VenuePackage, error) {
	var out entity.VenuePackage
	err := r.db.Table("category_package").
		Where("id = ?", ID).
		First(&out).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("[GetPackageByID] err: %v", err))
		}
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetPackageByID] err: %v", err))
	}
	return &out, nil
}

func (r *repository) GetGalleriesByVenueIDs(ctx context.Context, IDs []int) (map[int][]string, error) {
	var out []*VenueGallery
	err := r.db.Table("venue_gallery").
		Where("venue_id IN (?)", IDs).Find(&out).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetGalleriesByVenueIDs] err: %v", err))
	}

	urlsMappedByVenueID := map[int][]string{}
	for _, o := range out {
		urlsMappedByVenueID[o.VenueID] = append(urlsMappedByVenueID[o.VenueID], o.FileURL)
	}
	return urlsMappedByVenueID, nil
}

func (r *repository) GetOrdersByDate(ctx context.Context, date time.Time) ([]*entity.Order, error) {
	var out []*entity.Order
	err := r.db.Table("order").
		Where("date = ?", date).
		Find(&out).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetOrdersByDate] err: %v", err))
	}
	return out, nil
}

func (r *repository) GetVenuePackageByQuery(ctx context.Context, param *entity.GetVenuePackageQuery) ([]*entity.VenuePackage, error) {
	var dto []*VenuePackage
	qb := r.db.Table("category_package")
	if len(param.IDs) > 0 {
		qb = qb.Where("id IN (?)", param.IDs)
	}

	if len(param.CategoryIDs) > 0 {
		qb = qb.Where("category_id IN (?)", param.CategoryIDs)
	}

	err := qb.Find(&dto).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetVenuePackageByQuery] err: %v", err))
	}

	out := []*entity.VenuePackage{}
	for _, dt := range dto {
		out = append(out, dt.ToEntity())
	}
	return out, nil
}

func (r *repository) GetVenueCategoryPackageByQuery(ctx context.Context, param *entity.GetVenueCategoryByQuery) ([]*entity.VenuePackageCategory, error) {
	var dto []*VenuePackageCategory
	qb := r.db.Table("venue_category_package")
	if len(param.IDs) > 0 {
		qb = qb.Where("id IN (?)", param.IDs)
	}
	if param.VenueID > 0 {
		qb = qb.Where("venue_id = ?", param.VenueID)
	}

	err := qb.Find(&dto).Error
	if err != nil {
		return nil, errutil.New(errutil.ErrGeneralDB, fmt.Errorf("[GetVenuePackageByQuery] err: %v", err))
	}

	out := []*entity.VenuePackageCategory{}
	for _, dt := range dto {
		out = append(out, dt.ToEntity())
	}
	return out, nil
}
