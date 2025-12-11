package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID        int64
	Timestamp time.Time
	User      string
	Action    string
	Target    string
	Status    string
	Details   string
}

// AuditService handles audit logging
type AuditService struct {
	db *sql.DB
	mu sync.Mutex
}

var auditInstance *AuditService
var auditOnce sync.Once

// GetAuditService returns the singleton audit service
func GetAuditService() (*AuditService, error) {
	var initErr error
	auditOnce.Do(func() {
		auditInstance, initErr = newAuditService()
	})
	if initErr != nil {
		return nil, initErr
	}
	return auditInstance, nil
}

func newAuditService() (*AuditService, error) {
	// Use /data directory for persistent storage in K8s
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		// Fallback to current directory
		dataDir = "."
	}

	dbPath := filepath.Join(dataDir, "audit.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit database: %w", err)
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			user TEXT NOT NULL,
			action TEXT NOT NULL,
			target TEXT NOT NULL,
			status TEXT NOT NULL,
			details TEXT
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create audit table: %w", err)
	}

	return &AuditService{db: db}, nil
}

// Log records an audit entry
func (s *AuditService) Log(user, action, target, status, details string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		"INSERT INTO audit_log (user, action, target, status, details) VALUES (?, ?, ?, ?, ?)",
		user, action, target, status, details,
	)
	return err
}

// GetEntries retrieves audit entries, most recent first
func (s *AuditService) GetEntries(limit int) ([]AuditEntry, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.Query(`
		SELECT id, timestamp, user, action, target, status, COALESCE(details, '') 
		FROM audit_log 
		ORDER BY timestamp DESC 
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		var ts string
		if err := rows.Scan(&e.ID, &ts, &e.User, &e.Action, &e.Target, &e.Status, &e.Details); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse("2006-01-02 15:04:05", ts)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// Close closes the database connection
func (s *AuditService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
