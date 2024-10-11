package v1

import (
	"auth-service/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(router *gin.RouterGroup) {
	user := router.Group("")
	{
		user.POST("/sign-up", h.signUp)
		user.POST("/sign-in", h.signIn)
		user.POST("")
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var requestSignUp model.UserSignUp
	if err := c.BindJSON(&requestSignUp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.services.Users.SignUp(c.Request.Context(), requestSignUp)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.SetCookie("session", session.Token, session.ExpiresAt.Second(), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *Handler) signIn(c *gin.Context) {
	var userSignIn model.UserSignIn
	if err := c.BindJSON(&userSignIn); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.services.Users.SignIn(c.Request.Context(), userSignIn)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.SetCookie("session", session.Token, session.ExpiresAt.Second(), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User signed in successfully"})
}
