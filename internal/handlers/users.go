package handlers

import (
"fmt"
"html"
"log"
"net/http"
"strings"

"github.com/go-chi/chi/v5"
)

// ListUsers renders the users page for a domain
func ListUsers(w http.ResponseWriter, r *http.Request) {
domain := chi.URLParam(r, "domain")

w.Header().Set("Content-Type", "text/html")
w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Users - %s - MailHub Admin</title>
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
        <nav><a href="/domains">← Back to Domains</a></nav>
        <h1>Users: %s</h1>
        <button hx-get="/domains/%s/users/new" hx-target="#modal" hx-swap="innerHTML">Add User</button>
        <div id="user-list" hx-get="/domains/%s/users/list" hx-trigger="load" hx-swap="innerHTML">
            Loading...
        </div>
        <div id="modal"></div>
    </main>
</body>
</html>`, 
html.EscapeString(domain),
html.EscapeString(domain),
html.EscapeString(domain),
html.EscapeString(domain))))
}

// ListUsersPartial returns user list as HTML partial (for HTMX)
func ListUsersPartial(w http.ResponseWriter, r *http.Request) {
domain := chi.URLParam(r, "domain")
w.Header().Set("Content-Type", "text/html")

if h == nil || h.Mail == nil {
w.Write([]byte(`<p class="error">Mail service not initialized</p>`))
return
}

users, err := h.Mail.ListMailboxes(domain)
if err != nil {
log.Printf("Error listing users for %s: %v", domain, err)
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
return
}

if len(users) == 0 {
w.Write([]byte(`
<table>
    <thead>
        <tr>
            <th>Email</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td colspan="2"><em>No users configured yet</em></td>
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
            <th>Email</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>`)

for _, u := range users {
sb.WriteString(fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>
                <button class="secondary outline" 
                        hx-get="/domains/%s/users/%s/edit" 
                        hx-target="#modal" 
                        hx-swap="innerHTML">Change Password</button>
                <button class="secondary outline" 
                        hx-delete="/domains/%s/users/%s" 
                        hx-target="#user-list" 
                        hx-swap="innerHTML"
                        hx-confirm="Delete user %s?">Delete</button>
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
<dialog open>
    <article>
        <h3>Add User to %s</h3>
        <form hx-post="/domains/%s/users" hx-target="#user-list" hx-swap="innerHTML">
            <label>
                Username (before @)
                <input type="text" name="username" placeholder="user" required 
                       pattern="[a-zA-Z0-9._-]+">
            </label>
            <label>
                Password
                <input type="password" name="password" required minlength="8">
            </label>
            <footer>
                <button type="button" class="secondary" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Add User</button>
            </footer>
        </form>
    </article>
</dialog>`,
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
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
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
<dialog open>
    <article>
        <h3>Change Password: %s@%s</h3>
        <form hx-put="/domains/%s/users/%s/password" hx-target="#user-list" hx-swap="innerHTML">
            <label>
                New Password
                <input type="password" name="password" required minlength="8">
            </label>
            <footer>
                <button type="button" class="secondary" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Change Password</button>
            </footer>
        </form>
    </article>
</dialog>`,
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
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
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
w.Write([]byte(fmt.Sprintf(`<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))))
return
}

log.Printf("User deleted: %s@%s", user, domain)

// Return updated list
ListUsersPartial(w, r)
}
