package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chattcp/chattcp-capture/api"
	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line arguments
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	router := gin.Default()
	// CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// API Route
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/interfaces", api.ListInterfaces)
		apiGroup.GET("/capture", api.StartCaptureSSE)
	}
	addr := ":" + *port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
