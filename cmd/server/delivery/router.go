package delivery

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"code-playground/ui"
)

func NewRouter(rateLimit int, handler *SnippetHandler) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())
	r.Use(rateLimitMiddleware(rateLimit))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/run", handler.RunSnippet)
		v1.POST("/format", handler.FormatSnippet)
		v1.POST("/snippet/:id", handler.GetSnippet)
		v1.DELETE("/snippet/:id", handler.DeleteSnippet)
		v1.GET("/languages", handler.GetLanguages)
	}

	r.NoRoute(uiHandler())

	return r
}

func rateLimitMiddleware(limit int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen int64
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(c *gin.Context) {
		if limit <= 0 {
			c.Next()
			return
		}

		ip := c.ClientIP()
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(limit), limit)}
		}
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
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

		if _, err := ui.Static.Open(path[1:]); err == nil {
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

		// Inject dynamic meta tags for language landing pages
		lang := strings.Trim(path, "/")
		if info, ok := getLanguageMeta(lang); ok {
			html := string(index)
			
			// Replace Titles
			html = strings.Replace(html, "<title>PolyRun - Online Code Runner & Playground (Go, Python, JS)</title>", "<title>"+info.Title+"</title>", 1)
			html = strings.ReplaceAll(html, "content=\"PolyRun - Online Code Runner & Playground\"", "content=\""+info.Title+"\"")
			html = strings.ReplaceAll(html, "content=\"PolyRun - Online Code Runner\"", "content=\""+info.Title+"\"")
			
			// Replace Descriptions
			if info.Description != "" {
				html = strings.Replace(html, "content=\"PolyRun is a secure, minimalist online code runner and playground. Write, compile, and execute Go, Python, and JavaScript code instantly in your browser. Self-hosted and docker-isolated.\"", "content=\""+info.Description+"\"", 1)
				html = strings.ReplaceAll(html, "content=\"Secure, minimalist online code runner. Execute Go, Python, and JavaScript instantly.\"", "content=\""+info.Description+"\"")
				html = strings.Replace(html, "\"description\": \"A secure, minimalist online code runner for Go, Python, and JavaScript.\"", "\"description\": \""+info.Description+"\"", 1)
			}
			
			// Update URLs
			fullURL := "https://polyrun.kaichenl.com/" + strings.ToLower(lang)
			html = strings.Replace(html, "href=\"https://polyrun.kaichenl.com\"", "href=\""+fullURL+"\"", 1)
			html = strings.ReplaceAll(html, "content=\"https://polyrun.kaichenl.com/\"", "content=\""+fullURL+"\"")

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	}
}

type languageMeta struct {
	Title       string
	Description string
}

func getLanguageMeta(lang string) (languageMeta, bool) {
	metas := map[string]languageMeta{
		"python": {
			Title:       "PolyRun - Online Python Code Runner & Playground",
			Description: "Run Python code online with PolyRun. A secure and fast Python playground with support for multiple versions.",
		},
		"go": {
			Title:       "PolyRun - Online Go Code Runner & Playground",
			Description: "Run Go code online with PolyRun. A secure and fast Golang playground with support for multiple versions.",
		},
		"golang": {
			Title:       "PolyRun - Online Go Code Runner & Playground",
			Description: "Run Go code online with PolyRun. A secure and fast Golang playground with support for multiple versions.",
		},
		"javascript": {
			Title:       "PolyRun - Online JavaScript Code Runner & Playground",
			Description: "Run JavaScript code online with PolyRun. A secure and fast JS playground for web developers.",
		},
		"rust": {
			Title:       "PolyRun - Online Rust Code Runner & Playground",
			Description: "Run Rust code online with PolyRun. A secure and fast Rust playground for systems programming.",
		},
		"cpp": {
			Title:       "PolyRun - Online C++ Code Runner & Playground",
			Description: "Run C++ code online with PolyRun. A secure and fast C++ compiler and playground.",
		},
	}

	info, ok := metas[strings.ToLower(lang)]
	return info, ok
}
