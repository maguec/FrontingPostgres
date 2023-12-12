package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"fmt"
	"net/http"
	"strconv"
)

func Lookup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not conver id to integer"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var profile Profile
	db.Where(&Profile{SecondaryId: fmt.Sprintf("user%d", id)}).First(&profile)

	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("profile for user%d not found", id),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}
