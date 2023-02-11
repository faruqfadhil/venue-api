package api

import (
	"errors"
	"net/http"

	errutil "github.com/faruqfadhil/venue-api/pkg/error"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data interface{}   `json:"data"`
	Meta *ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	Status       string `json:"status"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	Page         int    `json:"page,omitempty"`
	TotalPage    int    `json:"totalPage,omitempty"`
	CurrentItems int    `json:"currentItems,omitempty"`
	TotalItems   int    `json:"totalItems,omitempty"`
}

func ResponseSuccess(c *gin.Context, out interface{}, meta *ResponseMeta) {
	c.JSON(http.StatusOK, Response{
		Data: out,
		Meta: meta,
	})
}

func ResponseFailed(c *gin.Context, err error) {
	resp := internalServerErr(err)
	typeErr := errutil.GetTypeErr(err)
	if errors.Is(typeErr, errutil.ErrGeneralBadRequest) {
		resp = badRequestErr(err)
	}
	if errors.Is(typeErr, errutil.ErrGeneralNotFound) {
		resp = notFoundErr(err)
	}
	c.JSON(resp.Meta.Code, resp)
}

func internalServerErr(err error) *Response {
	return &Response{
		Meta: &ResponseMeta{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func notFoundErr(err error) *Response {
	return &Response{
		Meta: &ResponseMeta{
			Status:  "error",
			Code:    http.StatusNotFound,
			Message: err.Error(),
		},
	}
}

func badRequestErr(err error) *Response {
	return &Response{
		Meta: &ResponseMeta{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		},
	}
}
