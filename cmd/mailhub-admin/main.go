package main

import (
"log"
"net/http"
"os"

"github.com/Ingasti/mailhub-admin/internal/config"
"github.com/Ingasti/mailhub-admin/internal/handlers"
"github.com/Ingasti/mailhub-admin/internal/middleware"
"github.com/Ingasti/mailhub-admin/internal/services"
"github.com/go-chi/chi/v5"
chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
// Load configuration
cfg := config.Load()

	// Initialize SSH client
	sshClient := services.NewSSHClient(services.SSHConfig{
		Host:        cfg.SSH.Host,
		Port:        cfg.SSH.Port,
		User:        cfg.SSH.User,
		KeyPath:     cfg.SSH.KeyPath,
		JumpHost:    cfg.SSH.JumpHost,
		JumpUser:    cfg.SSH.JumpUser,
		JumpKeyPath: cfg.SSH.JumpKeyPath,
	})// Initialize mail service
mailService := services.NewMailService(sshClient)

// Test connection
if err := mailService.TestConnection(); err != nil {
log.Printf("WARNING: SSH connection test failed: %v", err)
log.Printf("Mail management features may not work until SSH is configured")
} else {
log.Printf("SSH connection to mail server established")
}

// Initialize handlers with dependencies
handlers.Init(mailService)

// Setup router
r := chi.NewRouter()

// Global middleware
r.Use(chimiddleware.Logger)
r.Use(chimiddleware.Recoverer)
r.Use(chimiddleware.RealIP)

// Health check (no auth required)
r.Get("/health", handlers.HealthCheck)

// Protected routes
r.Group(func(r chi.Router) {
r.Use(middleware.Auth(cfg))

// Dashboard
r.Get("/", handlers.Dashboard)

// Rspamd antispam dashboard page
r.Get("/rspamd", handlers.HandleRspamdDashboard)

// Domain management
r.Route("/domains", func(r chi.Router) {
r.Get("/", handlers.ListDomains)
r.Get("/list", handlers.ListDomainsPartial)
r.Get("/new", handlers.NewDomainForm)
r.Post("/", handlers.CreateDomain)
r.Delete("/{domain}", handlers.DeleteDomain)

// Users per domain
r.Route("/{domain}/users", func(r chi.Router) {
r.Get("/", handlers.ListUsers)
r.Get("/list", handlers.ListUsersPartial)
r.Get("/new", handlers.NewUserForm)
r.Post("/", handlers.CreateUser)
r.Get("/{user}/edit", handlers.EditUserForm)
r.Put("/{user}/password", handlers.ChangePassword)
r.Delete("/{user}", handlers.DeleteUser)
})
})

// Audit log
r.Get("/audit", handlers.AuditLog)
r.Get("/audit/entries", handlers.AuditEntriesPartial)

// Rspamd antispam management
r.Route("/rspamd", func(r chi.Router) {
r.Get("/status", handlers.HandleRspamdStatus)
r.Get("/metrics", handlers.HandleRspamdMetrics)
r.Get("/config", handlers.HandleRspamdConfig)
r.Put("/config", handlers.HandleRspamdConfigUpdate)
r.Post("/config", handlers.HandleRspamdConfigUpdate)
r.Get("/whitelist", handlers.HandleRspamdWhitelist)
r.Post("/whitelist", handlers.HandleRspamdWhitelistAdd)
r.Delete("/whitelist", handlers.HandleRspamdWhitelistRemove)
r.Get("/logs", handlers.HandleRspamdLogs)
r.Post("/service/start", handlers.HandleRspamdServiceStart)
r.Post("/service/stop", handlers.HandleRspamdServiceStop)
r.Post("/service/restart", handlers.HandleRspamdServiceRestart)
r.Get("/export", handlers.HandleRspamdExport)
})
})

// Static files
fileServer := http.FileServer(http.Dir("web/static"))
r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

// Start server
addr := ":" + cfg.Port
log.Printf("Starting MailHub Admin on %s", addr)
log.Printf("SSH Host: %s, User: %s", cfg.SSH.Host, cfg.SSH.User)
log.Printf("Dev mode: %v", cfg.DevMode)

if err := http.ListenAndServe(addr, r); err != nil {
log.Fatal(err)
os.Exit(1)
}
}
