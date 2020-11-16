package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yesseneon/todo"
)

func (h *Handler) register(c *gin.Context) {
	var user todo.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) login(c *gin.Context) {

}
