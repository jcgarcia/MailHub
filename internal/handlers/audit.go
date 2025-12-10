package handlers

import (
	"net/http"
)

// AuditLog renders the audit log page
func AuditLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Audit Log - MailHub Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>
<body>
    <main class="container">
        <nav><a href="/">← Back to Dashboard</a></nav>
        <h1>Audit Log</h1>
        <div id="audit-list" hx-get="/audit/entries" hx-trigger="load" hx-swap="innerHTML">
            Loading...
        </div>
    </main>
</body>
</html>`))
}

// AuditEntriesPartial returns audit entries as HTML partial (for HTMX)
func AuditEntriesPartial(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch audit entries from SQLite
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<table>
    <thead>
        <tr>
            <th>Timestamp</th>
            <th>User</th>
            <th>Action</th>
            <th>Target</th>
            <th>Status</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td colspan="5"><em>No audit entries yet</em></td>
        </tr>
    </tbody>
</table>`))
}
