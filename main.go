package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var db *sql.DB

// var tmpl *template.Template

type Product struct {
	ID          string
	Name        string
	Description string
}

func initDB() {
	var err error
	connStr := "user=alimchik dbname=toysdb password=lego1234 host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Delete all existing data in the table before creating a new one (optional)
	_, err = db.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal("Failed to delete data:", err)
	}

	// Create the table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database initialized successfully!")
}

func homeHandler(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description); err != nil {
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}
		products = append(products, product)
	}

	// Serve the index.html file from the static folder
	c.HTML(200, "index.html", gin.H{
		"products": products,
	})
}

func productHandler(c *gin.Context) {
	id := c.Param("id")

	var product Product
	err := db.QueryRow("SELECT id, name, description FROM products WHERE id=$1", id).
		Scan(&product.ID, &product.Name, &product.Description)

	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	// Serve the product.html file
	c.HTML(200, "product.html", gin.H{
		"product": product,
	})
}

func createHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Serve the create.html file
		c.HTML(200, "create.html", nil)
		return
	}

	if c.Request.Method == "POST" {
		id := c.DefaultPostForm("id", "")
		name := c.DefaultPostForm("name", "")
		description := c.DefaultPostForm("description", "")

		_, err := db.Exec("INSERT INTO products (id, name, description) VALUES ($1, $2, $3)", id, name, description)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save product"})
			return
		}

		c.Redirect(303, "/")
	}
}

func main() {
	initDB()
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()
	// Serve static files like CSS, JS, images from the "static" folder
	r.Static("/static", "./static")
	// Load HTML templates
	r.LoadHTMLGlob("template/*.html")

	// Routes
	r.GET("/", homeHandler)
	r.GET("/product/:id", productHandler)
	r.GET("/create", createHandler)

	r.GET("/sitemap.xml", func(c *gin.Context) {
		sitemap := `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
		<url>
			<loc>https://lego.xyz/product/toy-car</loc>
			<lastmod>2025-02-10</lastmod>
			<changefreq>daily</changefreq>
			<priority>0.8</priority>
		</url>
		<url>
			<loc>https://lego.xyz/product/toy-train</loc>
			<lastmod>2025-02-09</lastmod>
			<changefreq>daily</changefreq>
			<priority>0.8</priority>
		</url>
		<!-- Add more URLs here -->
	</urlset>`

		c.Header("Content-Type", "application/xml")
		c.String(200, sitemap)
	})

	r.POST("/create", createHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
