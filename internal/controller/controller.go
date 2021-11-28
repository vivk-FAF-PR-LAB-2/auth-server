package controller

import (
	"github.com/gin-gonic/gin"
	"inter-protocol-auth-server/internal/auth/handler"
	"inter-protocol-auth-server/pkg/auth/authorization"
)

type IController interface {
	RegisterRoutes(r *gin.Engine)
}

type controller struct {
	autorizer authorization.IAuthorization
}

func NewController(authorizer authorization.IAuthorization) IController {
	return &controller{
		autorizer: authorizer,
	}
}

func (c *controller) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/auth")
	handler.RegisterHTTPEndpoints(api, c.autorizer)
}
