package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maguec/redisfrontingpostgres/api"
	"github.com/redis/rueidis"
	"gorm.io/gorm"
)

func APIMiddleWare(redisConn rueidis.Client, dbconn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redisConn)
		c.Set("db", dbconn)
		c.Next()
	}
}

func main() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		panic(err)
	}

	dbconn := api.DbConn("localhost", 5432, "postgres", "PgDbFTW15", "profiles")

	router := gin.New()
	var setup string
	setup = "initial"

	router.Use(APIMiddleWare(client, dbconn))

	router.PATCH("/config/:setup", func(c *gin.Context) {
		setup = c.Param("setup")
		c.JSON(http.StatusOK, gin.H{
			"setup": setup,
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"setup": setup,
		})
	})

	router.GET("/profile/:id", func(c *gin.Context) {
		c.Set("setup", setup)
		c.Next()
		api.Lookup(c)

	})

	router.POST("/load", api.Dataload)
	router.Run()
}
