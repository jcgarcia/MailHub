package handlers

import (
	"encoding/json"
	"net/http"
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
	// TODO: Render dashboard template
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>MailHub Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>
<body>
    <main class="container">
        <h1>MailHub Admin</h1>
        <p>Mail server administration tool</p>
        <nav>
            <ul>
                <li><a href="/domains">Domains</a></li>
                <li><a href="/audit">Audit Log</a></li>
            </ul>
        </nav>
    </main>
</body>
</html>`))
}
