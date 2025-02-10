package handlers

import (
	"toyswebsite/database"
	"toyswebsite/models"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description); err != nil {
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}
		products = append(products, product)
	}

	c.HTML(200, "index.html", gin.H{
		"products": products,
	})
}
