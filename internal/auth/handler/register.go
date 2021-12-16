package handler

import (
	"github.com/gin-gonic/gin"
	"inter-protocol-auth-server/pkg/auth/authorization"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, autorizer authorization.IAuthorization) {
	h := newHandler(autorizer)

	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)

	router.GET("/confirm/:token", h.confirm)
}
