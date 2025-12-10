package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListDomains renders the domains page
func ListDomains(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with SSH service
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Domains - MailHub Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>
<body>
    <main class="container">
        <nav><a href="/">← Back to Dashboard</a></nav>
        <h1>Mail Domains</h1>
        <button hx-get="/domains/new" hx-target="#modal" hx-swap="innerHTML">Add Domain</button>
        <div id="domain-list" hx-get="/domains/list" hx-trigger="load" hx-swap="innerHTML">
            Loading...
        </div>
        <div id="modal"></div>
    </main>
</body>
</html>`))
}

// ListDomainsPartial returns domain list as HTML partial (for HTMX)
func ListDomainsPartial(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch domains via SSH and render list
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<table>
    <thead>
        <tr>
            <th>Domain</th>
            <th>Users</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td colspan="3"><em>No domains configured yet</em></td>
        </tr>
    </tbody>
</table>`))
}

// NewDomainForm returns the add domain form
func NewDomainForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<dialog open>
    <article>
        <h3>Add Domain</h3>
        <form hx-post="/domains" hx-target="#domain-list" hx-swap="innerHTML">
            <label>
                Domain Name
                <input type="text" name="domain" placeholder="example.com" required>
            </label>
            <footer>
                <button type="button" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Add Domain</button>
            </footer>
        </form>
    </article>
</dialog>`))
}

// CreateDomain adds a new mail domain
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with SSH service
	domain := r.FormValue("domain")
	if domain == "" {
		http.Error(w, "Domain required", http.StatusBadRequest)
		return
	}

	// TODO: Execute SSH command to add domain
	// TODO: Log to audit

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

	// TODO: Execute SSH command to delete domain
	// TODO: Log to audit

	w.WriteHeader(http.StatusOK)
}
