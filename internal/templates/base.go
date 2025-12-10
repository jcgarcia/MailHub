package templates

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// BaseCSS contains the shared CSS for all pages
const BaseCSS = `
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: linear-gradient(135deg, #e0e5ec 0%, #d0d5dc 100%);
    min-height: 100vh;
    padding: 20px;
}
.container {
    max-width: 900px;
    margin: 0 auto;
}
.card {
    background: white;
    border-radius: 16px;
    padding: 30px;
    box-shadow: 0 10px 40px rgba(0,0,0,0.1);
    margin-bottom: 20px;
}
.header {
    text-align: center;
    margin-bottom: 25px;
}
.header .logo {
    width: 60px;
    height: 60px;
    margin-bottom: 10px;
}
.header h1 {
    color: #1a73e8;
    font-size: 1.5rem;
    margin-bottom: 5px;
}
.header .subtitle {
    color: #666;
    font-size: 0.95rem;
}
.nav-link {
    display: inline-flex;
    align-items: center;
    color: #1a73e8;
    text-decoration: none;
    font-size: 0.9rem;
    margin-bottom: 15px;
}
.nav-link:hover {
    text-decoration: underline;
}
.nav-link i {
    margin-right: 6px;
}
.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 500;
    text-decoration: none;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
}
.btn-primary {
    background: #1a73e8;
    color: white;
}
.btn-primary:hover {
    background: #1557b0;
    transform: translateY(-1px);
}
.btn-secondary {
    background: #f8f9fa;
    color: #333;
    border: 1px solid #ddd;
}
.btn-secondary:hover {
    background: #e8f0fe;
    border-color: #1a73e8;
}
.btn-danger {
    background: #dc3545;
    color: white;
}
.btn-danger:hover {
    background: #c82333;
}
.btn-sm {
    padding: 8px 16px;
    font-size: 0.85rem;
}
table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
}
th, td {
    padding: 14px 16px;
    text-align: left;
    border-bottom: 1px solid #eee;
}
th {
    background: #f8f9fa;
    font-weight: 600;
    color: #333;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}
tr:hover {
    background: #f8f9fa;
}
.menu-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 15px;
    margin-top: 20px;
}
.menu-item {
    display: flex;
    align-items: center;
    padding: 18px 20px;
    background: #f8f9fa;
    border-radius: 10px;
    text-decoration: none;
    color: #333;
    font-size: 1rem;
    transition: all 0.2s ease;
    border: 1px solid #eee;
}
.menu-item:hover {
    background: #e8f0fe;
    border-color: #1a73e8;
    transform: translateX(4px);
}
.menu-item i {
    font-size: 1.4rem;
    margin-right: 14px;
    color: #1a73e8;
    width: 28px;
    text-align: center;
}
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0,0,0,0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}
.modal {
    background: white;
    border-radius: 16px;
    padding: 30px;
    width: 100%;
    max-width: 400px;
    box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.modal h3 {
    margin-bottom: 20px;
    color: #333;
}
.form-group {
    margin-bottom: 20px;
}
.form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #333;
}
.form-group input {
    width: 100%;
    padding: 12px 14px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.2s;
}
.form-group input:focus {
    outline: none;
    border-color: #1a73e8;
    box-shadow: 0 0 0 3px rgba(26,115,232,0.1);
}
.modal-footer {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    margin-top: 25px;
}
.error-msg {
    background: #fee2e2;
    color: #dc2626;
    padding: 12px 16px;
    border-radius: 8px;
    margin-bottom: 15px;
}
.success-msg {
    background: #dcfce7;
    color: #16a34a;
    padding: 12px 16px;
    border-radius: 8px;
    margin-bottom: 15px;
}
.empty-state {
    text-align: center;
    padding: 40px 20px;
    color: #666;
}
.empty-state i {
    font-size: 3rem;
    color: #ccc;
    margin-bottom: 15px;
}
.badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 0.8rem;
    font-weight: 500;
}
.badge-info {
    background: #e8f0fe;
    color: #1a73e8;
}
.badge-success {
    background: #dcfce7;
    color: #16a34a;
}
.actions {
    display: flex;
    gap: 8px;
}
`

// BaseHead returns the common HTML head section
func BaseHead(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - MailHub Admin</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/line-awesome/1.3.0/line-awesome/css/line-awesome.min.css">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>%s</style>
</head>
<body>`, template.HTMLEscapeString(title), BaseCSS)
}

// BaseFooter returns the common HTML footer
func BaseFooter() string {
	return `</body></html>`
}

// Logo returns the SVG logo
func Logo() string {
	return `<svg class="logo" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="50" cy="50" r="45" fill="#1a73e8"/>
    <path d="M25 35 L50 55 L75 35 L75 70 L25 70 Z" fill="white"/>
    <path d="M25 35 L50 55 L75 35" stroke="white" stroke-width="3" fill="none"/>
</svg>`
}

// RenderPage renders a full page with the base template
func RenderPage(w io.Writer, title string, content string) {
	var sb strings.Builder
	sb.WriteString(BaseHead(title))
	sb.WriteString(`<div class="container">`)
	sb.WriteString(content)
	sb.WriteString(`</div>`)
	sb.WriteString(BaseFooter())
	w.Write([]byte(sb.String()))
}
