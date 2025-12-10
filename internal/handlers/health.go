package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Ingasti/mailhub-admin/internal/templates"
)

// HealthCheck returns server health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"app":    "mailhub-admin",
	})
}

// Dashboard renders the main dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	content := `
<div class="card">
    <div class="header">
        ` + templates.Logo() + `
        <h1>MailHub Admin</h1>
        <p class="subtitle">Mail Server Administration</p>
    </div>
    
    <div class="menu-grid">
        <a href="/domains" class="menu-item">
            <i class="la la-globe"></i>
            <span>Domains</span>
        </a>
        <a href="/audit" class="menu-item">
            <i class="la la-history"></i>
            <span>Audit Log</span>
        </a>
    </div>
</div>`

	templates.RenderPage(w, "Dashboard", content)
}
