package handlers

import (
	"github.com/Ingasti/mailhub-admin/internal/services"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	Mail *services.MailService
}

// Global handler instance (initialized in main)
var h *Handler

// Init initializes the handler with dependencies
func Init(mail *services.MailService) {
	h = &Handler{
		Mail: mail,
	}
}

// GetHandler returns the handler instance
func GetHandler() *Handler {
	return h
}
