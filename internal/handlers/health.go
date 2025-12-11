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
</div>

<div class="card">
    <h2 style="color: #1a73e8; margin-bottom: 20px; font-size: 1.2rem;">
        <i class="la la-cog" style="margin-right: 8px;"></i>Mail Client Configuration
    </h2>
    <p style="color: #666; margin-bottom: 20px;">Use these settings to configure your email client (Outlook, Thunderbird, Apple Mail, etc.):</p>
    
    <div class="config-grid">
        <div class="config-section">
            <h3><i class="la la-envelope"></i> Incoming Mail (IMAP)</h3>
            <table class="config-table">
                <tr><td>Server:</td><td><strong>cmh.ingasti.com</strong></td></tr>
                <tr><td>Port:</td><td><strong>993</strong> (SSL/TLS)</td></tr>
                <tr><td>Security:</td><td>SSL/TLS</td></tr>
                <tr><td>Username:</td><td>Your full email address</td></tr>
            </table>
        </div>
        
        <div class="config-section">
            <h3><i class="la la-paper-plane"></i> Outgoing Mail (SMTP)</h3>
            <table class="config-table">
                <tr><td>Server:</td><td><strong>cmh.ingasti.com</strong></td></tr>
                <tr><td>Port:</td><td><strong>587</strong> (STARTTLS)</td></tr>
                <tr><td>Security:</td><td>STARTTLS</td></tr>
                <tr><td>Username:</td><td>Your full email address</td></tr>
                <tr><td>Auth:</td><td>Password authentication</td></tr>
            </table>
        </div>
    </div>
    
    <div style="background: #e8f0fe; padding: 15px; border-radius: 8px; margin-top: 20px;">
        <p style="margin: 0; color: #1a73e8; font-size: 0.9rem;">
            <i class="la la-info-circle" style="margin-right: 6px;"></i>
            <strong>Note:</strong> Use your full email address (e.g., user@yourdomain.com) as both username and email address.
        </p>
    </div>
</div>`

	templates.RenderPage(w, "Dashboard", content)
}
