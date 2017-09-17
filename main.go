package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	handlers "github.com/tjw0051/log-go/Handlers"
	store "github.com/tjw0051/log-go/Store"
)

var (
	masterAPIKey = os.Getenv("MASTER_KEY")
	dbHostname   = os.Getenv("DB_HOSTNAME")
	dbPort       = os.Getenv("DB_PORT")
	dbUsername   = os.Getenv("DB_USERNAME")
	dbPassword   = os.Getenv("DB_PASSWORD")
	dbName       = os.Getenv("DB_NAME")
)

// https://docs.docker.com/compose/networking/

func main() {

	// Connect to Database
	err := store.Connect(dbHostname, dbPort, dbUsername, dbName, dbPassword)
	if err != nil {
		panic("FATAL ERROR: Could not connect to DB! \n " + err.Error())
	} else {
		log.Println("Connected to database.")
	}

	// Register - define healthcheck route
	// server polls healthcheck route

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")
	{
		// API Key Management
		keys := v1.Group("/keys")
		keys.Use(MasterAuthMiddleware())
		{
			keys.POST("", handlers.CreateKeys)
			keys.GET("", handlers.GetKeys)
			keys.DELETE("", handlers.DeleteKeys)
		}

		// Logging
		log := v1.Group("/log")
		log.Use(KeyAuthMiddleware())
		{
			log.POST("", handlers.CreateLog)
		}

		query := v1.Group("/query")
		query.Use(MasterAuthMiddleware())
		{
			query.GET("", handlers.Query)
		}
	}

	router.Run()
}

var loadedKeys = []string{}

func KeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")

		// Bad auth header
		splitToken := strings.Split(token, " ")
		if len(splitToken) != 2 {
			c.Status(401)
			c.Abort()
			return
		}

		// Wrong auth type
		if strings.ToLower(splitToken[0]) != "basic" {
			c.Status(401)
			c.Abort()
			return
		}

		// Find key in cache
		for i := 0; i < len(loadedKeys); i++ {
			if loadedKeys[i] == splitToken[1] {
				c.Next()
				return
			}
		}

		// Reload cache
		keys, _ := store.GetKeys()
		for i := 0; i < len(keys); i++ {
			loadedKeys = append(loadedKeys, keys[i].Key)
		}

		// Find key in cache (again)
		for i := 0; i < len(loadedKeys); i++ {
			if loadedKeys[i] == splitToken[1] {
				c.Next()
				return
			}
		}

		c.Status(401)
		c.Abort()
	}
}

func MasterAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")

		// Bad auth header
		splitToken := strings.Split(token, " ")
		if len(splitToken) != 2 {
			c.Status(401)
			c.Abort()
			return
		}

		// Wrong auth type
		if strings.ToLower(splitToken[0]) != "basic" {
			c.Status(401)
			c.Abort()
			return
		}

		// Wrong key
		if splitToken[1] != masterAPIKey {
			c.Status(401)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORS Support
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // http://192.168.0.22:3000
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
