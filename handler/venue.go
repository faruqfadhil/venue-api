package handler

import (
	"context"
	"net/http"

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
