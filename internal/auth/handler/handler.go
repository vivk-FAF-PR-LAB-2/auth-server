package handler

import (
	"github.com/gin-gonic/gin"
	authError "inter-protocol-auth-server/internal/auth/error"
	"inter-protocol-auth-server/internal/models"
	"inter-protocol-auth-server/pkg/auth/authorization"
	"net/http"
)

type handler struct {
	autorizer authorization.IAuthorization
}

func newHandler(autorizer authorization.IAuthorization) *handler {
	return &handler{
		autorizer: autorizer,
	}
}

func (h *handler) signUp(c *gin.Context) {
	inp := new(models.User)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	if err := h.autorizer.SignUp(c.Request.Context(), inp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "confirm user email"))
}

func (h *handler) confirm(c *gin.Context) {
	token := c.Param("token")
	if err := h.autorizer.Confirm(c.Request.Context(), token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "user created successfully"))
}

func (h *handler) signIn(c *gin.Context) {
	inp := new(models.User)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.autorizer.SignIn(c.Request.Context(), inp)
	if err != nil {
		if err == authError.ErrInvalidAccessToken {
			c.AbortWithStatusJSON(http.StatusBadRequest, newSignInResponse(STATUS_ERROR, err.Error(), ""))
			return
		}

		if err == authError.ErrUserDoesNotExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, newSignInResponse(STATUS_ERROR, err.Error(), ""))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, newSignInResponse(STATUS_ERROR, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, newSignInResponse(STATUS_OK, "", token))
}
