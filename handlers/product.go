package handlers

import (
	"toyswebsite/database"
	"toyswebsite/models"

	"github.com/gin-gonic/gin"
)

func ProductHandler(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	err := database.DB.QueryRow("SELECT id, name, description FROM products WHERE id=$1", id).
		Scan(&product.ID, &product.Name, &product.Description)

	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	c.HTML(200, "product.html", gin.H{
		"product": product,
	})
}
