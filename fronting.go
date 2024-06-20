package main

import (
	"fmt"
	"net/http"

	"gitlab.com/avarf/getenvs"

	"github.com/gin-gonic/gin"
	"github.com/maguec/redisfrontingpostgres/api"
	"github.com/redis/rueidis"
	"gorm.io/gorm"
)

func APIMiddleWare(redisConn rueidis.Client, dbconn *gorm.DB, datasize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redisConn)
		c.Set("db", dbconn)
		c.Set("datasize", datasize)
		c.Next()
	}
}

func main() {
	dbserver := getenvs.GetEnvString("PGHOST", "localhost")
	dbport, _ := getenvs.GetEnvInt("PGPORT", 5432)
	datasize, _ := getenvs.GetEnvInt("DATASIZE", 100000)
	dbuser := getenvs.GetEnvString("PGUSER", "postgres")
	dbpassword := getenvs.GetEnvString("PGPASSWORD", "PgDbFTW15")
	dbname := getenvs.GetEnvString("PGDB", "profiles")
	redisserver := getenvs.GetEnvString("REDIS_SERVER", "localhost")
	redisport, _ := getenvs.GetEnvInt("REDIS_PORT", 6379)
	rediscache, _ := getenvs.GetEnvBool( "REDIS_CACHE", false,) // By default we do not use redis client side caching

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", redisserver, redisport)}, // []string{"localhost:6379"},
    DisableCache: !rediscache,
	})
	if err != nil {
		panic(err)
	}

	dbconn := api.DbConn(dbserver, dbport, dbuser, dbpassword, dbname)

	router := gin.New()
	var setup string
	setup = "initial"

	router.Use(APIMiddleWare(client, dbconn, datasize))

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
