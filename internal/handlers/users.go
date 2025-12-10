package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListUsers renders the users page for a domain
func ListUsers(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Users - ` + domain + ` - MailHub Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>
<body>
    <main class="container">
        <nav><a href="/domains">← Back to Domains</a></nav>
        <h1>Users: ` + domain + `</h1>
        <button hx-get="/domains/` + domain + `/users/new" hx-target="#modal" hx-swap="innerHTML">Add User</button>
        <div id="user-list" hx-get="/domains/` + domain + `/users/list" hx-trigger="load" hx-swap="innerHTML">
            Loading...
        </div>
        <div id="modal"></div>
    </main>
</body>
</html>`))
}

// ListUsersPartial returns user list as HTML partial (for HTMX)
func ListUsersPartial(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch users via SSH and render list
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<table>
    <thead>
        <tr>
            <th>Email</th>
            <th>Created</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td colspan="3"><em>No users configured yet</em></td>
        </tr>
    </tbody>
</table>`))
}

// NewUserForm returns the add user form
func NewUserForm(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<dialog open>
    <article>
        <h3>Add User</h3>
        <form hx-post="/domains/` + domain + `/users" hx-target="#user-list" hx-swap="innerHTML">
            <label>
                Username (before @)
                <input type="text" name="username" placeholder="user" required>
            </label>
            <label>
                Password
                <input type="password" name="password" required minlength="8">
            </label>
            <footer>
                <button type="button" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Add User</button>
            </footer>
        </form>
    </article>
</dialog>`))
}

// CreateUser adds a new email user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" || domain == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	// TODO: Execute SSH command to add user
	// TODO: Log to audit

	// Return updated list
	ListUsersPartial(w, r)
}

// EditUserForm returns the edit user form
func EditUserForm(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	user := chi.URLParam(r, "user")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<dialog open>
    <article>
        <h3>Change Password: ` + user + `@` + domain + `</h3>
        <form hx-put="/domains/` + domain + `/users/` + user + `/password" hx-target="#user-list" hx-swap="innerHTML">
            <label>
                New Password
                <input type="password" name="password" required minlength="8">
            </label>
            <footer>
                <button type="button" onclick="this.closest('dialog').remove()">Cancel</button>
                <button type="submit">Change Password</button>
            </footer>
        </form>
    </article>
</dialog>`))
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

	// TODO: Execute SSH command to change password
	// TODO: Log to audit

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

	// TODO: Execute SSH command to delete user
	// TODO: Log to audit

	w.WriteHeader(http.StatusOK)
}
