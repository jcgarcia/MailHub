package handlers

import (
	"fmt"
	"net/http"

	"github.com/Ingasti/mailhub-admin/internal/services"
	"github.com/Ingasti/mailhub-admin/internal/templates"
)

// AuditLog renders the audit log page
func AuditLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	content := `
<div class="card">
    <a href="/" class="nav-link"><i class="la la-arrow-left"></i> Back to Dashboard</a>
    <div class="header">
        ` + templates.Logo() + `
        <h1>Audit Log</h1>
        <p class="subtitle">Activity History</p>
    </div>
    
    <div id="audit-list" hx-get="/audit/entries" hx-trigger="load" hx-swap="innerHTML">
        <p style="text-align: center; color: #666;">Loading audit entries...</p>
    </div>
</div>`

	templates.RenderPage(w, "Audit Log", content)
}

// AuditEntriesPartial returns audit entries as HTML partial (for HTMX)
func AuditEntriesPartial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	auditSvc, err := services.GetAuditService()
	if err != nil {
		w.Write([]byte(`<div class="alert alert-danger">Failed to load audit service</div>`))
		return
	}

	entries, err := auditSvc.GetEntries(50)
	if err != nil {
		w.Write([]byte(`<div class="alert alert-danger">Failed to load audit entries</div>`))
		return
	}

	if len(entries) == 0 {
		w.Write([]byte(`
<div style="text-align: center; padding: 40px; color: #666;">
    <i class="la la-clipboard-list" style="font-size: 3rem; margin-bottom: 15px; display: block; opacity: 0.5;"></i>
    <p><em>No audit entries yet</em></p>
    <p style="font-size: 0.9rem;">Activity will be logged when you manage domains and mailboxes.</p>
</div>`))
		return
	}

	html := `<table>
    <thead>
        <tr>
            <th>Timestamp</th>
            <th>User</th>
            <th>Action</th>
            <th>Target</th>
            <th>Status</th>
        </tr>
    </thead>
    <tbody>`

	for _, e := range entries {
		statusClass := "badge-success"
		if e.Status == "failed" {
			statusClass = "badge-danger"
		}
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td><span class="badge %s">%s</span></td>
        </tr>`,
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.User,
			e.Action,
			e.Target,
			statusClass,
			e.Status,
		)
	}

	html += `
    </tbody>
</table>`

	w.Write([]byte(html))
}

// LogAudit is a helper to log audit entries from handlers
func LogAudit(user, action, target, status, details string) {
	auditSvc, err := services.GetAuditService()
	if err != nil {
		return // Silently fail - audit logging shouldn't break the app
	}
	auditSvc.Log(user, action, target, status, details)
}
