package entity

import "time"

type Venue struct {
	ID           int
	Name         string
	MinPrice     float64
	MaxPrice     float64
	Capacity     int
	Star         float64
	ReviewCount  int
	ThumbnailURL string
	Gallery      []string
	City         *City
}

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VenueDetail struct {
	ID          int
	Name        string
	Description string
	Website     string
	Phone       string
	Email       string
	Instagram   string
	Address     string
	Logo        string
	Categories  []*VenuePackageCategory
}

type VenuePackageCategory struct {
	Description string
	Packages    []*VenuePackage
}

type VenuePackage struct {
	ID           int
	ThumbnailURL string
	Name         string
	Price        float64
	Capacity     int
}

type GetVenuesParam struct {
	CityID      int
	IsFavourite bool
	Date        time.Time
	Page        int
	Limit       int
}

type Pagination struct {
	Page         int
	TotalPage    int
	CurrentItems int
	TotalItems   int
}
