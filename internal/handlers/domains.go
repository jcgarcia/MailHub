package handlers

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/Ingasti/mailhub-admin/internal/templates"
	"github.com/go-chi/chi/v5"
)

// ListDomains renders the domains page
func ListDomains(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	content := `
<div class="card">
    <a href="/" class="nav-link"><i class="la la-arrow-left"></i> Back to Dashboard</a>
    <div class="header">
        <h1>Mail Domains</h1>
        <p class="subtitle">Manage your mail domains</p>
    </div>
    
    <button class="btn btn-primary" hx-get="/domains/new" hx-target="#modal" hx-swap="innerHTML">
        <i class="la la-plus" style="margin-right: 8px;"></i> Add Domain
    </button>
    
    <div id="domain-list" hx-get="/domains/list" hx-trigger="load" hx-swap="innerHTML">
        <div class="empty-state">
            <i class="la la-spinner la-spin"></i>
            <p>Loading domains...</p>
        </div>
    </div>
    <div id="modal"></div>
</div>`

	templates.RenderPage(w, "Domains", content)
}

// ListDomainsPartial returns domain list as HTML partial (for HTMX)
func ListDomainsPartial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if h == nil || h.Mail == nil {
		w.Write([]byte(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Mail service not initialized</div>`))
		return
	}

	domains, err := h.Mail.ListDomains()
	if err != nil {
		log.Printf("Error listing domains: %v", err)
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	if len(domains) == 0 {
		w.Write([]byte(`
<div class="empty-state">
    <i class="la la-inbox"></i>
    <p>No domains configured yet</p>
    <p style="font-size: 0.9rem; margin-top: 10px;">Click "Add Domain" to get started</p>
</div>`))
		return
	}

	var sb strings.Builder
	sb.WriteString(`
<table>
    <thead>
        <tr>
            <th>Domain</th>
            <th>Users</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>`)

	for _, d := range domains {
		sb.WriteString(fmt.Sprintf(`
        <tr>
            <td>
                <a href="/domains/%s/users" style="color: #1a73e8; text-decoration: none; font-weight: 500;">
                    <i class="la la-globe" style="margin-right: 6px;"></i>%s
                </a>
            </td>
            <td><span class="badge badge-info">%d users</span></td>
            <td class="actions">
                <button class="btn btn-danger btn-sm" 
                        hx-delete="/domains/%s" 
                        hx-target="#domain-list" 
                        hx-swap="innerHTML"
                        hx-confirm="Delete domain %s and all its users?">
                    <i class="la la-trash"></i>
                </button>
            </td>
        </tr>`,
			html.EscapeString(d.Name),
			html.EscapeString(d.Name),
			d.UserCount,
			html.EscapeString(d.Name),
			html.EscapeString(d.Name)))
	}

	sb.WriteString(`
    </tbody>
</table>`)

	w.Write([]byte(sb.String()))
}

// NewDomainForm returns the add domain form
func NewDomainForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<div class="modal-overlay" onclick="if(event.target===this) this.remove()">
    <div class="modal">
        <h3><i class="la la-plus-circle" style="color: #1a73e8; margin-right: 8px;"></i>Add Domain</h3>
        <form hx-post="/domains" hx-target="#domain-list" hx-swap="innerHTML">
            <div class="form-group">
                <label for="domain">Domain Name</label>
                <input type="text" id="domain" name="domain" placeholder="example.com" required 
                       pattern="[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}">
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button type="submit" class="btn btn-primary">Add Domain</button>
            </div>
        </form>
    </div>
</div>`))
}

// CreateDomain adds a new mail domain
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.FormValue("domain"))
	domain = strings.ToLower(domain)

	if domain == "" {
		http.Error(w, "Domain required", http.StatusBadRequest)
		return
	}

	if h == nil || h.Mail == nil {
		http.Error(w, "Mail service not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.Mail.AddDomain(domain); err != nil {
		log.Printf("Error adding domain %s: %v", domain, err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	log.Printf("Domain added: %s", domain)

	// Return updated list
	ListDomainsPartial(w, r)
}

// DeleteDomain removes a mail domain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	if domain == "" {
		http.Error(w, "Domain required", http.StatusBadRequest)
		return
	}

	if h == nil || h.Mail == nil {
		http.Error(w, "Mail service not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.Mail.DeleteDomain(domain); err != nil {
		log.Printf("Error deleting domain %s: %v", domain, err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	log.Printf("Domain deleted: %s", domain)

	// Return updated list
	ListDomainsPartial(w, r)
}
