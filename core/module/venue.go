package module

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	Order(ctx context.Context, order *entity.Order) error
	GetVenuesNearby(ctx context.Context) ([]*entity.VenueNearby, error)
	GetVenueByID(ctx context.Context, ID int) (*entity.VenueDetail, error)
	GetPackageByID(ctx context.Context, ID int) (*entity.PackageDetail, error)
}

type usecase struct {
	repo repository.Repository
}

func New(repo repository.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetVenues(ctx context.Context, param entity.GetVenuesParam) ([]*entity.Venue, *entity.Pagination, error) {
	if !param.Date.IsZero() {
		param.Date = time.Date(param.Date.Year(), param.Date.Month(), param.Date.Day(), 0, 0, 0, 0, time.UTC)
		existingOrders, err := u.repo.GetOrdersByDate(ctx, param.Date)
		if err != nil {
			return nil, nil, err
		}

		packageIDs := []int{}
		for _, eo := range existingOrders {
			packageIDs = append(packageIDs, eo.PackageID)
		}

		categoryIDs := []int{}
		if len(packageIDs) > 0 {
			packages, err := u.repo.GetVenuePackageByQuery(ctx, &entity.GetVenuePackageQuery{
				IDs: packageIDs,
			})
			if err != nil {
				return nil, nil, err
			}
			for _, pkg := range packages {
				categoryIDs = append(categoryIDs, pkg.CategoryID)
			}
		}

		venueIDs := []int{}
		if len(categoryIDs) > 0 {
			categories, err := u.repo.GetVenueCategoryPackageByQuery(ctx, &entity.GetVenueCategoryByQuery{
				IDs: categoryIDs,
			})
			if err != nil {
				return nil, nil, err
			}
			for _, ctg := range categories {
				venueIDs = append(venueIDs, ctg.VenueID)
			}
		}

		if len(venueIDs) > 0 {
			param.NotInIDs = venueIDs
		}
	}

	venues, pag, err := u.repo.GetVenues(ctx, param)
	if err != nil {
		return nil, nil, err
	}
	venueIDs := []int{}
	venuesMappedByCityID := map[int][]*entity.Venue{}
	for _, vn := range venues {
		venueIDs = append(venueIDs, vn.ID)
		venuesMappedByCityID[vn.CityID] = append(venuesMappedByCityID[vn.CityID], vn)
	}

	// Map city
	if len(venuesMappedByCityID) > 0 {
		cities, err := u.repo.GetCities(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, ct := range cities {
			if vns, ok := venuesMappedByCityID[ct.ID]; ok {
				for _, vn := range vns {
					vn.City = ct
				}
			}
		}
	}

	// Map gallery
	if len(venueIDs) > 0 {
		galleryMappedByVenueID, err := u.repo.GetGalleriesByVenueIDs(ctx, venueIDs)
		if err != nil {
			return nil, nil, err
		}
		for _, vn := range venues {
			if gl, ok := galleryMappedByVenueID[vn.ID]; ok {
				vn.Gallery = gl
			}
		}
	}
	return venues, pag, nil
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

func (u *usecase) Order(ctx context.Context, order *entity.Order) error {
	_, err := u.repo.GetPackageByID(ctx, order.PackageID)
	if err != nil {
		if errors.Is(errutil.GetTypeErr(err), errutil.ErrGeneralNotFound) {
			return errutil.New(errutil.ErrGeneralBadRequest, err, fmt.Sprintf("Tidak dapat membuat order untuk tanggal %v dikarenakan package id %d tidak ditemukan", order.Date, order.PackageID))
		}
		return err
	}

	order.Date = time.Date(order.Date.Year(), order.Date.Month(), order.Date.Day(), 0, 0, 0, 0, time.UTC)
	existingOrder, err := u.repo.GetOrderByPackageIDAndDate(ctx, order.PackageID, order.Date)
	if err != nil && !errors.Is(errutil.GetTypeErr(err), errutil.ErrGeneralNotFound) {
		return err
	}
	if existingOrder != nil {
		return errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("unavailable date"), fmt.Sprintf("Tidak dapat membuat order untuk tanggal %v dikarenakan tempat sudah di reservasi", order.Date))
	}

	return u.repo.CreateOrder(ctx, order)
}

func (u *usecase) GetVenuesNearby(ctx context.Context) ([]*entity.VenueNearby, error) {
	cities, err := u.repo.GetCities(ctx)
	if err != nil {
		return nil, err
	}
	cityIDs := []int{}
	cityMappedByID := map[int]*entity.City{}
	for _, ct := range cities {
		cityIDs = append(cityIDs, ct.ID)
		cityMappedByID[ct.ID] = ct
	}

	venues, _, err := u.repo.GetVenues(ctx, entity.GetVenuesParam{
		CityIDs:             cityIDs,
		IsWithoutPagination: true,
	})
	if err != nil {
		return nil, err
	}

	venuesMappedByCityId := map[int][]*entity.Venue{}
	for _, vn := range venues {
		venuesMappedByCityId[vn.CityID] = append(venuesMappedByCityId[vn.CityID], vn)
	}

	out := []*entity.VenueNearby{}
	for cityId, vns := range venuesMappedByCityId {
		if _, ok := cityMappedByID[cityId]; !ok {
			continue
		}
		if len(vns) > 0 {
			out = append(out, &entity.VenueNearby{
				CityID:       cityId,
				CityName:     cityMappedByID[cityId].Name,
				TotalVenue:   len(vns),
				ThumbnailURL: vns[0].ThumbnailURL,
			})
		}
	}

	return out, nil
}

func (u *usecase) GetVenueByID(ctx context.Context, ID int) (*entity.VenueDetail, error) {
	venues, _, err := u.GetVenues(ctx, entity.GetVenuesParam{
		ID:                  ID,
		IsWithoutPagination: true,
	})
	if err != nil {
		return nil, err
	}

	if len(venues) < 1 {
		return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("venue not found"), "venue tidak ditemukan")
	}

	categories, err := u.repo.GetVenueCategoryPackageByQuery(ctx, &entity.GetVenueCategoryByQuery{
		VenueID: ID,
	})
	if err != nil {
		return nil, err
	}

	categoryIDs := []int{}
	for _, ctg := range categories {
		categoryIDs = append(categoryIDs, ctg.ID)
	}

	if len(categoryIDs) > 0 {
		packages, err := u.repo.GetVenuePackageByQuery(ctx, &entity.GetVenuePackageQuery{
			CategoryIDs: categoryIDs,
		})
		if err != nil {
			return nil, err
		}
		packagesMappedByCategoryID := map[int][]*entity.VenuePackage{}
		for _, pkg := range packages {
			packagesMappedByCategoryID[pkg.CategoryID] = append(packagesMappedByCategoryID[pkg.CategoryID], pkg)
		}

		for _, ctg := range categories {
			if pkgs, ok := packagesMappedByCategoryID[ctg.ID]; ok {
				ctg.Packages = pkgs
			}
		}
	}

	return &entity.VenueDetail{
		ID:          venues[0].ID,
		Name:        venues[0].Name,
		Description: venues[0].Description,
		Website:     venues[0].Website,
		Phone:       venues[0].Phone,
		Email:       venues[0].Email,
		Instagram:   venues[0].Instagram,
		Address:     venues[0].Address,
		Logo:        venues[0].Logo,
		Categories:  categories,
	}, nil
}

func (u *usecase) GetPackageByID(ctx context.Context, ID int) (*entity.PackageDetail, error) {
	pkg, err := u.repo.GetVenuePackageByQuery(ctx, &entity.GetVenuePackageQuery{
		IDs: []int{ID},
	})
	if err != nil {
		return nil, err
	}
	if len(pkg) < 1 {
		return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("package not found"), "package tidak ditemukan")
	}
	category, err := u.repo.GetVenueCategoryPackageByQuery(ctx, &entity.GetVenueCategoryByQuery{
		IDs: []int{pkg[0].CategoryID},
	})
	if err != nil {
		return nil, err
	}
	if len(category) < 1 {
		return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("category not found"), "category tidak ditemukan")
	}
	venue, _, err := u.GetVenues(ctx, entity.GetVenuesParam{
		ID: category[0].VenueID,
	})
	if err != nil {
		return nil, err
	}
	if len(venue) < 1 {
		return nil, errutil.New(errutil.ErrGeneralNotFound, fmt.Errorf("venue not found"), "venue tidak ditemukan")
	}

	return &entity.PackageDetail{
		ID:           pkg[0].ID,
		ThumbnailURL: pkg[0].ThumbnailURL,
		Name:         pkg[0].Name,
		Price:        pkg[0].Price,
		Capacity:     pkg[0].Capacity,
		VenueName:    venue[0].Name,
		VenuePhone:   venue[0].Phone,
		Gallery:      venue[0].Gallery,
		Description:  pkg[0].Description,
	}, nil
}
