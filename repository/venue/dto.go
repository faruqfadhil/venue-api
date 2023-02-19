package venue

import "github.com/faruqfadhil/venue-api/core/entity"

type Venue struct {
	ID           int
	Name         string
	MinPrice     float64
	MaxPrice     float64
	Capacity     int
	Star         float64
	ReviewCount  int
	ThumbnailURL string
	Description  string
	Website      string
	Phone        string
	Email        string
	Instagram    string
	Address      string
	Logo         string
	IsFavourite  bool
	CityID       int
}

func (v *Venue) ToEntity() *entity.Venue {
	return &entity.Venue{
		ID:           v.ID,
		Name:         v.Name,
		MinPrice:     v.MinPrice,
		MaxPrice:     v.MaxPrice,
		Capacity:     v.Capacity,
		Star:         v.Star,
		ReviewCount:  v.ReviewCount,
		ThumbnailURL: v.ThumbnailURL,
		Description:  v.Description,
		Website:      v.Website,
		Phone:        v.Phone,
		Email:        v.Email,
		Instagram:    v.Instagram,
		Address:      v.Address,
		Logo:         v.Logo,
		IsFavourite:  v.IsFavourite,
		CityID:       v.CityID,
	}
}

type VenueGallery struct {
	ID      int
	VenueID int
	FileURL string
}

type VenuePackageCategory struct {
	ID          int
	VenueID     int
	Description string
}

func (v *VenuePackageCategory) ToEntity() *entity.VenuePackageCategory {
	return &entity.VenuePackageCategory{
		ID:          v.ID,
		VenueID:     v.VenueID,
		Description: v.Description,
	}
}

type VenuePackage struct {
	ID           int
	CategoryID   int
	ThumbnailURL string
	Name         string
	Price        float64
	Capacity     int
	Description  string
}

func (v *VenuePackage) ToEntity() *entity.VenuePackage {
	return &entity.VenuePackage{
		ID:           v.ID,
		CategoryID:   v.CategoryID,
		ThumbnailURL: v.ThumbnailURL,
		Name:         v.Name,
		Price:        v.Price,
		Capacity:     v.Capacity,
		Description:  v.Description,
	}
}
