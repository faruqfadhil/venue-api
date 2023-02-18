package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/faruqfadhil/venue-api/core/module"
	errutil "github.com/faruqfadhil/venue-api/pkg/error"
	"github.com/gin-gonic/gin"
)

type MiddlewareService struct {
	authSvc module.Usecase
}

func NewMiddlewareService(authSvc module.Usecase) *MiddlewareService {
	return &MiddlewareService{
		authSvc: authSvc,
	}
}

func (s *MiddlewareService) AuthenticateRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if !strings.Contains(token, "Bearer") {
			ResponseFailed(ctx, errutil.New(errutil.ErrUnauthorized, fmt.Errorf("invalid token"), "anda tidak diizinkan mengakses aplikasi ini"))
			ctx.Abort()
			return
		}
		token = strings.Replace(token, "Bearer ", "", -1)

		validate, err := s.authSvc.ValidateToken(context.Background(), token)
		if err != nil {
			ResponseFailed(ctx, errutil.New(errutil.ErrUnauthorized, err, "anda tidak diizinkan mengakses aplikasi ini"))
			ctx.Abort()
			return
		}
		if validate == nil {
			ResponseFailed(ctx, errutil.New(errutil.ErrUnauthorized, err, "anda tidak diizinkan mengakses aplikasi ini"))
			ctx.Abort()
			return
		}

		if validate != nil {
			ctx.Set("id", validate.ID)
			ctx.Set("email", validate.Email)
			ctx.Set("fullname", validate.FullName)
			ctx.Next()
			return
		}
	}
}
