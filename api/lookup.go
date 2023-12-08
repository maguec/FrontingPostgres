package api

import (
	"github.com/gin-gonic/gin"
	//"github.com/redis/rueidis"
	"net/http"
	"strconv"
)

func Lookup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not conver id to integer"})
	} else {
	c.JSON(http.StatusOK, gin.H{
		"setup": c.MustGet("setup").(string),
		"id":    id,
	})
}
}

func Dataload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "loading data",
	})
}
