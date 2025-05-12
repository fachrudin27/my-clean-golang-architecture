package gingo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGin() *gin.Engine {
	var app = gin.New()

	app.Use(NotfoundMiddleware())
	app.Use(NewErrorHandler())
	app.Use(CORSMiddleware())
	app.Use(gin.Logger(), gin.Recovery())
	app.SetTrustedProxies([]string{"127.0.0.1"})
	return app
}

func NewErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NotfoundMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Continue to the next middleware or route handler

		// Check if no route was matched and status is 404
		if c.Writer.Status() == http.StatusNotFound {
			// Send custom 404 response
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Sorry. We couldn't find that page",
			})
			c.Abort() // Stop further processing
			return
		}
	}
}
