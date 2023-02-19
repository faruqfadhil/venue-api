package entity

import (
	"time"
)

type Venue struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	MinPrice     float64  `json:"minPrice"`
	MaxPrice     float64  `json:"maxPrice"`
	Capacity     int      `json:"capacity"`
	Star         float64  `json:"star"`
	ReviewCount  int      `json:"reviewCount"`
	ThumbnailURL string   `json:"thumbnailUrl"`
	CityID       int      `json:"cityID"`
	City         *City    `json:"city"`
	Description  string   `json:"description"`
	Website      string   `json:"website"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Instagram    string   `json:"instagram"`
	Address      string   `json:"address"`
	Logo         string   `json:"logo"`
	IsFavourite  bool     `json:"isFavourite"`
	Gallery      []string `json:"gallery"`
}

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VenueDetail struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Website     string                  `json:"website"`
	Phone       string                  `json:"phone"`
	Email       string                  `json:"email"`
	Instagram   string                  `json:"instagram"`
	Address     string                  `json:"address"`
	Logo        string                  `json:"logo"`
	Categories  []*VenuePackageCategory `json:"categories"`
}

type VenuePackageCategory struct {
	ID          int             `json:"id"`
	VenueID     int             `json:"venueId"`
	Description string          `json:"description"`
	Packages    []*VenuePackage `json:"packages"`
}

type VenuePackage struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"categoryId"`
	ThumbnailURL string  `json:"thumbnailUrl"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Capacity     int     `json:"capacity"`
	Description  string  `json:"description"`
}

type GetVenuesParam struct {
	ID                  int
	CityID              int
	CityIDs             []int
	IsFavourite         bool
	Date                time.Time
	Page                int
	Limit               int
	NotInIDs            []int
	IsWithoutPagination bool
}

type Pagination struct {
	Page         int
	TotalPage    int
	CurrentItems int
	TotalItems   int
}

type Order struct {
	ID        int
	PackageID int
	UserID    int
	Date      time.Time
}

type GetVenuePackageQuery struct {
	IDs         []int
	CategoryIDs []int
}

type GetVenueCategoryByQuery struct {
	IDs     []int
	VenueID int
}

type VenueNearby struct {
	CityID       int    `json:"cityId"`
	CityName     string `json:"cityName"`
	TotalVenue   int    `json:"totalVenue"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type PackageDetail struct {
	ID           int      `json:"id"`
	ThumbnailURL string   `json:"thumbnailUrl"`
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	Capacity     int      `json:"capacity"`
	VenueName    string   `json:"venueName"`
	VenuePhone   string   `json:"venuePhone"`
	Gallery      []string `json:"gallery"`
	Description  string   `json:"description"`
}
