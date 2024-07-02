package api

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/hints"

	"fmt"
	"net/http"
	"strconv"
	"time"
)

func Lookup(c *gin.Context) {
	var profile Profile
	var redisLoad bool
	id, err := strconv.Atoi(c.Param("id"))
	logger := c.MustGet("logger").(*zap.SugaredLogger)
	start := time.Now()
	if err != nil {
		logger.Errorw("lookup", "error", err, "elapsed", time.Since(start).Milliseconds(), "id", id)
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
				logger.Errorw("unmarshal", "error", err, "elapsed", time.Since(start).Milliseconds(), "id", id)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not decode profile from redis"})
				return
			}
			logger.Infow("lookup", "source", "redis", "elapsed", time.Since(start).String(), "id", id)
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
		logger.Errorw("dberror", "error", err, "elapsed", time.Since(start).Milliseconds(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile.ID == 0 {
		logger.Errorw("profileerror", "error", "notFound", "elapsed", time.Since(start).Milliseconds(), "id", id)
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
			logger.Errorw("redis-set-error", "error", "notFound", "elapsed", time.Since(start).Milliseconds(), "id", id)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	logger.Infow("lookup", "source", "database", "elapsed", time.Since(start).String(), "id", id)
	c.JSON(http.StatusOK, profile)
}
