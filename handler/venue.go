package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
