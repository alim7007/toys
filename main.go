package main

import (
	"log"

	"toyswebsite/database"
	"toyswebsite/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDB()
	defer database.DB.Close()

	// Initialize Gin router
	r := gin.Default()
	r.Static("/template", "./template")
	r.LoadHTMLGlob("template/*.html")

	// Routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/product/:id", handlers.ProductHandler)
	r.GET("/create", handlers.CreateHandler)
	r.POST("/create", handlers.CreateHandler)
	r.GET("/sitemap.xml", handlers.SitemapHandler)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
