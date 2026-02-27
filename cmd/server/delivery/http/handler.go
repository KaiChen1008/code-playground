package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"code-playground/cmd/server/domain"
	"code-playground/cmd/server/domain/models"
)

type SnippetHandler struct {
	uc domain.Usecase
}

func NewSnippetHandler(uc domain.Usecase) *SnippetHandler {
	return &SnippetHandler{uc: uc}
}

func (h *SnippetHandler) FormatSnippet(c *gin.Context) {
	var req models.FormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if possible, though Gin does some binding

	resp, err := h.uc.FormatSnippet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SnippetHandler) RunSnippet(c *gin.Context) {
	var req models.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.RunSnippet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SnippetHandler) GetSnippet(c *gin.Context) {
	id := c.Param("id")
	snippet, err := h.uc.GetSnippet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) DeleteSnippet(c *gin.Context) {
	id := c.Param("id")
	err := h.uc.DeleteSnippet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snippet not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SnippetHandler) GetLanguages(c *gin.Context) {
	languages, err := h.uc.GetLanguages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, languages)
}
