package handlers

import (
	"database/sql"
	"toyswebsite/database"
	"toyswebsite/models"

	"github.com/gin-gonic/gin"
)

func CreateHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Serve the create.html file
		c.HTML(200, "create.html", nil)
		return
	}

	if c.Request.Method == "POST" {
		// Parse form data
		id := c.DefaultPostForm("id", "")
		name := c.DefaultPostForm("name", "")
		description := c.DefaultPostForm("description", "")

		// Validate if the product ID already exists
		var existingProduct models.Product
		err := database.DB.QueryRow("SELECT id FROM products WHERE id = $1", id).Scan(&existingProduct.ID)
		if err != nil && err != sql.ErrNoRows {
			// Database error
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}

		if err == nil {
			// Product with the same ID already exists
			c.HTML(400, "create.html", gin.H{
				"error":       "Product ID already exists",
				"id":          id,
				"name":        name,
				"description": description,
			})
			return
		}

		// Insert the new product into the database
		_, err = database.DB.Exec("INSERT INTO products (id, name, description) VALUES ($1, $2, $3)", id, name, description)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save product"})
			return
		}

		// Redirect to the home page after successful creation
		c.Redirect(303, "/")
	}
}
