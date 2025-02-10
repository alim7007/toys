package handlers

import (
	"fmt"
	"toyswebsite/database"

	"github.com/gin-gonic/gin"
)

func SitemapHandler(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var sitemap string
	sitemap += `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}

		productURL := fmt.Sprintf("https://lego.xyz/product/%s", id)
		sitemap += fmt.Sprintf(`
		<url>
			<loc>%s</loc>
			<lastmod>%s</lastmod>
			<changefreq>daily</changefreq>
			<priority>0.8</priority>
		</url>`, productURL, "2025-02-10")
	}

	sitemap += `</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(200, sitemap)
}
