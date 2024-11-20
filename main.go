package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var client *firestore.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	ctx := context.Background()
	credentials := os.Getenv("FIREBASE_CREDENTIALS")

	client, err = firestore.NewClient(ctx, "gdgoc-backend-6ba98", option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		log.Fatalf("Failed to create Firstore client: %v", err)
	}
	defer client.Close()

	router := gin.Default()

	// Get array of books
	router.GET("/api/books", getBooks)
	// Not found
	router.GET("/api/books/:id", getBook)
	// Create new book
	router.POST("/api/books", createBook)
	// Update book
	router.PUT("/api/books/2", updateBook)
	// Delete book
	router.DELETE("/api/books/2", deleteBook)

	// Run server on port 8000
	router.Run(":8000")
}

func getBooks(c *gin.Context) {
	ctx := context.Background()
	var items []map[string]interface{}
	iter := client.Collection("gdgoc").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		items = append(items, doc.Data())
	}
	c.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func getBook(c *gin.Context) {
	response := gin.H{
		"message": "Book not found",
	}
	c.JSON(http.StatusNotFound, response)
}

func createBook(c *gin.Context) {
	ctx := context.Background()
	var item map[string]interface{}
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item["created_at"] = "2024-10-25T13:36:09.000000Z"
	item["updated_at"] = "2024-10-25T13:36:09.000000Z"
	item["id"] = 2
	_, err := client.Collection("gdgoc").Doc("2").Set(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message": "Book created successfully",
		"data": gin.H{
			"title":        item["title"],
			"author":       item["author"],
			"published_at": item["published_at"],
			"updated_at":   item["updated_at"],
			"created_at":   item["created_at"],
			"id":           item["id"],
		},
	}
	c.JSON(http.StatusCreated, response)
}

func updateBook(c *gin.Context) {
	ctx := context.Background()
	var item map[string]interface{}
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := client.Collection("gdgoc").Doc("2").Set(ctx, item, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doc, err := client.Collection("gdgoc").Doc("2").Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"data":    doc.Data(),
	})
}

func deleteBook(c *gin.Context) {
	ctx := context.Background()
	_, err := client.Collection("items").Doc("2").Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
