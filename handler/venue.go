package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	errutil "github.com/faruqfadhil/venue-api/pkg/error"

	"github.com/faruqfadhil/venue-api/core/entity"
	"github.com/faruqfadhil/venue-api/core/module"
	"github.com/faruqfadhil/venue-api/pkg/api"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	usecase module.Usecase
}

func New(uc module.Usecase) *HTTPHandler {
	return &HTTPHandler{
		usecase: uc,
	}
}

func (h *HTTPHandler) GetCities(c *gin.Context) {
	cities, err := h.usecase.GetCities(context.Background())
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, cities, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}

type HTTPRegister struct {
	Data *entity.User `json:"data"`
}

func (h *HTTPHandler) Register(c *gin.Context) {
	var payload *HTTPRegister
	if err := c.ShouldBindJSON(&payload); err != nil {
		api.ResponseFailed(c, errutil.ErrGeneralBadRequest)
		return
	}
	if strings.TrimSpace(payload.Data.Email) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("email can't be empty"), "email tidak boleh kosong"))
		return
	}
	if strings.TrimSpace(payload.Data.Password) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("password can't be empty"), "password tidak boleh kosong"))
		return
	}
	if strings.TrimSpace(payload.Data.FullName) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("fullname can't be empty"), "fullname tidak boleh kosong"))
		return
	}
	err := h.usecase.Register(context.Background(), payload.Data)
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, nil, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusCreated,
	})
}

type HTTPLogin struct {
	Data *HTTPLoginData `json:"data"`
}

type HTTPLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HTTPLoginResp struct {
	Account *entity.Auth `json:"account"`
}

func (h *HTTPHandler) Login(c *gin.Context) {
	var payload *HTTPLogin
	if err := c.ShouldBindJSON(&payload); err != nil {
		api.ResponseFailed(c, errutil.ErrGeneralBadRequest)
		return
	}
	if strings.TrimSpace(payload.Data.Email) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("email can't be empty"), "email tidak boleh kosong"))
		return
	}
	if strings.TrimSpace(payload.Data.Password) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("password can't be empty"), "password tidak boleh kosong"))
		return
	}
	authInfo, err := h.usecase.Login(context.Background(), payload.Data.Email, payload.Data.Password)
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	resp := HTTPLoginResp{
		Account: authInfo,
	}

	api.ResponseSuccess(c, resp, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}

type HTTPOrder struct {
	Data *HTTPOrderData `json:"data"`
}
type HTTPOrderData struct {
	PackageID int    `json:"packageId"`
	Date      string `json:"date"`
}

func (h *HTTPHandler) CreateOrder(c *gin.Context) {
	var payload *HTTPOrder
	if err := c.ShouldBindJSON(&payload); err != nil {
		api.ResponseFailed(c, errutil.ErrGeneralBadRequest)
		return
	}
	if strings.TrimSpace(payload.Data.Date) == "" {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("date can't be empty"), "tanggal tidak boleh kosong"))
		return
	}
	if payload.Data.PackageID < 1 {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("package id can't be empty"), "package id tidak boleh kosong"))
		return
	}
	date, err := time.Parse("2006-01-02", payload.Data.Date)
	if err != nil {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid date format"), "format tanggal harus YYYY-MM-DD"))
		return
	}
	if _, ok := c.Get("id"); !ok {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("can't extract user id"), "tidak dapat mengekstrak user id"))
		return
	}
	userID, _ := c.Get("id")
	err = h.usecase.Order(context.Background(), &entity.Order{
		PackageID: payload.Data.PackageID,
		UserID:    userID.(int),
		Date:      date,
	})
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}
	api.ResponseSuccess(c, nil, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}

type HTTPVenues struct {
	Venues []*entity.Venue `json:"venues"`
}

func (h *HTTPHandler) GetVenues(c *gin.Context) {
	var (
		cityID      int
		isFavourite bool
		date        time.Time
		page        int
		limit       int
	)
	cityIDQ := c.Query("cityId")
	if cityIDQ != "" {
		city, err := strconv.Atoi(cityIDQ)
		if err != nil {
			api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid date format"), "format city id tidak valid"))
			return
		}
		cityID = city
	}

	isF, ok := c.GetQuery("isFavourite")
	if ok && isF == "true" {
		isFavourite = true
	}
	d := c.Query("date")
	if d != "" {
		dn, err := time.Parse("2006-01-02", d)
		if err != nil {
			api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid date format"), "format tanggal harus YYYY-MM-DD"))
			return
		}
		date = dn
	}

	pageQ := c.Query("page")
	if pageQ != "" {
		t, err := strconv.Atoi(pageQ)
		if err != nil {
			api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid date format"), "format page tidak valid"))
			return
		}
		page = t
	}
	limitQ := c.Query("limit")
	if limitQ != "" {
		t, err := strconv.Atoi(limitQ)
		if err != nil {
			api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid date format"), "format limit tidak valid"))
			return
		}
		limit = t
	}

	result, pag, err := h.usecase.GetVenues(c, entity.GetVenuesParam{
		CityID:      cityID,
		IsFavourite: isFavourite,
		Date:        date,
		Page:        page,
		Limit:       limit,
	})
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, HTTPVenues{
		Venues: result,
	}, &api.ResponseMeta{
		Status:       "success",
		Code:         http.StatusOK,
		Page:         pag.Page,
		TotalPage:    pag.TotalPage,
		CurrentItems: pag.CurrentItems,
		TotalItems:   pag.TotalItems,
	})
}

type HTTPGetNearby struct {
	Nearbies []*entity.VenueNearby `json:"nearbies"`
}

func (h *HTTPHandler) GetNearby(c *gin.Context) {
	result, err := h.usecase.GetVenuesNearby(context.Background())
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, HTTPGetNearby{
		Nearbies: result,
	}, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}

type HTTPGetVenueDetail struct {
	Venue *entity.VenueDetail `json:"venue"`
}

func (h *HTTPHandler) GetVenueDetail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid id format"), "format id tidak valid"))
		return
	}

	result, err := h.usecase.GetVenueByID(c, idInt)
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, HTTPGetVenueDetail{
		Venue: result,
	}, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}

type HTTPGetPackageDetail struct {
	Package *entity.PackageDetail `json:"package"`
}

func (h *HTTPHandler) GetPackageDetail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		api.ResponseFailed(c, errutil.New(errutil.ErrGeneralBadRequest, fmt.Errorf("invalid id format"), "format id tidak valid"))
		return
	}

	result, err := h.usecase.GetPackageByID(c, idInt)
	if err != nil {
		api.ResponseFailed(c, err)
		return
	}

	api.ResponseSuccess(c, HTTPGetPackageDetail{
		Package: result,
	}, &api.ResponseMeta{
		Status: "success",
		Code:   http.StatusOK,
	})
}
