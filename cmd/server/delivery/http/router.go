package http

import (
	"code-playground/ui"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *SnippetHandler) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/run", handler.RunSnippet)
		v1.POST("/format", handler.FormatSnippet)
		v1.GET("/snippet/:id", handler.GetSnippet)
		v1.DELETE("/snippet/:id", handler.DeleteSnippet)
		v1.GET("/languages", handler.GetLanguages)
	}

	r.NoRoute(uiHandler())

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func uiHandler() gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(ui.Static))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/v1") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		if _, err := ui.Static.Open(path); err == nil {
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		if strings.Contains(path, ".") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		index, err := ui.Static.ReadFile("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load UI"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}
