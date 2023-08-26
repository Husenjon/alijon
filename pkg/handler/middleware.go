package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	cookie, err := c.Cookie("session")
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth token")
		return
	}
	userId, err := h.services.Authoration.ParseToken(cookie)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}
