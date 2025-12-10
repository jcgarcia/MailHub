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

// ListUsers renders the users page for a domain
func ListUsers(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	w.Header().Set("Content-Type", "text/html")

	content := fmt.Sprintf(`
<div class="card">
    <a href="/domains" class="nav-link"><i class="la la-arrow-left"></i> Back to Domains</a>
    <div class="header">
        <h1>Users: %s</h1>
        <p class="subtitle">Manage mailboxes for this domain</p>
    </div>
    
    <button class="btn btn-primary" hx-get="/domains/%s/users/new" hx-target="#modal" hx-swap="innerHTML">
        <i class="la la-user-plus" style="margin-right: 8px;"></i> Add User
    </button>
    
    <div id="user-list" hx-get="/domains/%s/users/list" hx-trigger="load" hx-swap="innerHTML">
        <div class="empty-state">
            <i class="la la-spinner la-spin"></i>
            <p>Loading users...</p>
        </div>
    </div>
    <div id="modal"></div>
</div>`,
		html.EscapeString(domain),
		html.EscapeString(domain),
		html.EscapeString(domain))

	templates.RenderPage(w, "Users - "+domain, content)
}

// ListUsersPartial returns user list as HTML partial (for HTMX)
func ListUsersPartial(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	w.Header().Set("Content-Type", "text/html")

	if h == nil || h.Mail == nil {
		w.Write([]byte(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Mail service not initialized</div>`))
		return
	}

	users, err := h.Mail.ListMailboxes(domain)
	if err != nil {
		log.Printf("Error listing users for %s: %v", domain, err)
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	if len(users) == 0 {
		w.Write([]byte(`
<div class="empty-state">
    <i class="la la-user-slash"></i>
    <p>No users configured yet</p>
    <p style="font-size: 0.9rem; margin-top: 10px;">Click "Add User" to create a mailbox</p>
</div>`))
		return
	}

	var sb strings.Builder
	sb.WriteString(`
<table>
    <thead>
        <tr>
            <th>Email</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>`)

	for _, u := range users {
		sb.WriteString(fmt.Sprintf(`
        <tr>
            <td>
                <i class="la la-envelope" style="color: #1a73e8; margin-right: 8px;"></i>
                <strong>%s</strong>
            </td>
            <td class="actions">
                <button class="btn btn-secondary btn-sm" 
                        hx-get="/domains/%s/users/%s/edit" 
                        hx-target="#modal" 
                        hx-swap="innerHTML">
                    <i class="la la-key"></i>
                </button>
                <button class="btn btn-danger btn-sm" 
                        hx-delete="/domains/%s/users/%s" 
                        hx-target="#user-list" 
                        hx-swap="innerHTML"
                        hx-confirm="Delete user %s?">
                    <i class="la la-trash"></i>
                </button>
            </td>
        </tr>`,
			html.EscapeString(u.Email),
			html.EscapeString(domain),
			html.EscapeString(u.Username),
			html.EscapeString(domain),
			html.EscapeString(u.Username),
			html.EscapeString(u.Email)))
	}

	sb.WriteString(`
    </tbody>
</table>`)

	w.Write([]byte(sb.String()))
}

// NewUserForm returns the add user form
func NewUserForm(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
<div class="modal-overlay" onclick="if(event.target===this) this.remove()">
    <div class="modal">
        <h3><i class="la la-user-plus" style="color: #1a73e8; margin-right: 8px;"></i>Add User to %s</h3>
        <form hx-post="/domains/%s/users" hx-target="#user-list" hx-swap="innerHTML">
            <div class="form-group">
                <label for="username">Username (before @)</label>
                <input type="text" id="username" name="username" placeholder="user" required 
                       pattern="[a-zA-Z0-9._-]+">
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required minlength="8">
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button type="submit" class="btn btn-primary">Add User</button>
            </div>
        </form>
    </div>
</div>`,
		html.EscapeString(domain),
		html.EscapeString(domain))))
}

// CreateUser adds a new email user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	username := strings.TrimSpace(r.FormValue("username"))
	username = strings.ToLower(username)
	password := r.FormValue("password")

	if username == "" || password == "" || domain == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if h == nil || h.Mail == nil {
		http.Error(w, "Mail service not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.Mail.AddMailbox(domain, username, password); err != nil {
		log.Printf("Error adding user %s@%s: %v", username, domain, err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	log.Printf("User added: %s@%s", username, domain)

	// Return updated list
	ListUsersPartial(w, r)
}

// EditUserForm returns the edit user form
func EditUserForm(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	user := chi.URLParam(r, "user")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
<div class="modal-overlay" onclick="if(event.target===this) this.remove()">
    <div class="modal">
        <h3><i class="la la-key" style="color: #1a73e8; margin-right: 8px;"></i>Change Password</h3>
        <p style="color: #666; margin-bottom: 20px;">%s@%s</p>
        <form hx-put="/domains/%s/users/%s/password" hx-target="#user-list" hx-swap="innerHTML">
            <div class="form-group">
                <label for="password">New Password</label>
                <input type="password" id="password" name="password" required minlength="8">
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button type="submit" class="btn btn-primary">Change Password</button>
            </div>
        </form>
    </div>
</div>`,
		html.EscapeString(user),
		html.EscapeString(domain),
		html.EscapeString(domain),
		html.EscapeString(user))))
}

// ChangePassword updates a user's password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	user := chi.URLParam(r, "user")
	password := r.FormValue("password")

	if password == "" || domain == "" || user == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	if h == nil || h.Mail == nil {
		http.Error(w, "Mail service not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.Mail.ChangePassword(domain, user, password); err != nil {
		log.Printf("Error changing password for %s@%s: %v", user, domain, err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	log.Printf("Password changed for: %s@%s", user, domain)

	// Return updated list
	ListUsersPartial(w, r)
}

// DeleteUser removes an email user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	user := chi.URLParam(r, "user")

	if domain == "" || user == "" {
		http.Error(w, "Domain and user required", http.StatusBadRequest)
		return
	}

	if h == nil || h.Mail == nil {
		http.Error(w, "Mail service not initialized", http.StatusInternalServerError)
		return
	}

	if err := h.Mail.DeleteMailbox(domain, user); err != nil {
		log.Printf("Error deleting user %s@%s: %v", user, domain, err)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<div class="error-msg"><i class="la la-exclamation-circle"></i> Error: %s</div>`, html.EscapeString(err.Error()))))
		return
	}

	log.Printf("User deleted: %s@%s", user, domain)

	// Return updated list
	ListUsersPartial(w, r)
}
