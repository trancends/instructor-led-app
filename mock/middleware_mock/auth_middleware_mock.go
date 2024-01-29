package middleware_mock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type AuthMiddlewareMock struct {
	mock.Mock
}

func (m *AuthMiddlewareMock) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
