package api

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
	"gorm.io/gorm"
	"gorm.io/hints"

	"fmt"
	"net/http"
	"strconv"
)

func Lookup(c *gin.Context) {
	var profile Profile
	var redisLoad bool
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not conver id to integer"})
		return
	}
	setup := c.MustGet("setup").(string)
	redis := c.MustGet("redis").(rueidis.Client)

	if setup == "caching" {
		kn := fmt.Sprintf("profile:%d", id)
		val, err := redis.Do(context.Background(), redis.B().Get().Key(kn).Build()).ToString()
		if err == nil && val != "" {
			// IF we find in Redis - return it
			err := json.Unmarshal([]byte(val), &profile)
			if err != nil {
				fmt.Printf("%+v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not decode profile from redis"})
				return
			}
			c.JSON(http.StatusOK, profile)
			// IF we find in Redis - return it
			return
		} else {
			// Load the record to Redis later
			redisLoad = true
		}
	}
	db := c.MustGet("db").(*gorm.DB)
	db.Clauses(hints.CommentAfter("limit", "route='/lookup',module='api.Lookup'")).Where(&Profile{SecondaryId: fmt.Sprintf("user%d", id)}).First(&profile)

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

	if redisLoad {
		val, _ := json.Marshal(&profile)
		kn := fmt.Sprintf("profile:%d", id)
		err = redis.Do(context.Background(), redis.B().Set().Key(kn).Value(string(val)).Build()).Error()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, profile)
}
