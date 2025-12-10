package handlers

import (
"fmt"
"html"
"log"
"net/http"
"strings"

"github.com/go-chi/chi/v5"
)

// ListDomains renders the domains page
func ListDomains(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Domains - MailHub Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
    <style>
        .error { color: var(--pico-del-color); }
        .success { color: var(--pico-ins-color); }
        dialog { max-width: 400px; }
    </style>
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
w.Header().Set("Content-Type", "text/html")

if h == nil || h.Mail == nil {
w.Write([]byte(`<p class="error">Mail service not initialized</p>`))
return
}

domains, err := h.Mail.ListDomains()
if err != nil {
log.Printf("Error listing domains: %v", err)
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
return
}

if len(domains) == 0 {
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
            <td><a href="/domains/%s/users">%s</a></td>
            <td>%d</td>
            <td>
                <button class="secondary outline" 
                        hx-delete="/domains/%s" 
                        hx-target="#domain-list" 
                        hx-swap="innerHTML"
                        hx-confirm="Delete domain %s and all its users?">Delete</button>
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
<dialog open>
    <article>
        <h3>Add Domain</h3>
        <form hx-post="/domains" hx-target="#domain-list" hx-swap="innerHTML">
            <label>
                Domain Name
                <input type="text" name="domain" placeholder="example.com" required 
                       pattern="[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}">
            </label>
            <footer>
                <button type="button" class="secondary" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Add Domain</button>
            </footer>
        </form>
    </article>
</dialog>`))
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
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
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
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
return
}

log.Printf("Domain deleted: %s", domain)

// Return updated list
ListDomainsPartial(w, r)
}
