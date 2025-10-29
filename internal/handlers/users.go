package handlers

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/db/daos"
	"github.com/Stenoliv/didlydoodash_api/internal/db/models"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Get all users paginated
func GetAllUsers(c *gin.Context) {
	result, err := daos.GetUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ServerError)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Get a single user
func GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = *models.CurrentUser
	}
	usr, err := daos.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ServerError)
		return
	}

	if (usr == models.User{}) { // Zero-value check
		c.JSON(http.StatusNotFound, utils.UserNotFound)
		return
	}

	c.JSON(http.StatusOK, usr)
}

// Put profile/user
func PutUser(c *gin.Context) {

	c.JSON(http.StatusOK, nil)
}

// Patch profile/user
func PatchUser(c *gin.Context) {

	c.JSON(http.StatusOK, nil)
}
