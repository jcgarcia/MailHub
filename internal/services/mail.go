package services

import (
	"fmt"
	"sort"
	"strings"
)

// MailService provides mail server management operations
type MailService struct {
	ssh *SSHClient
}

// Domain represents a mail domain
type Domain struct {
	Name      string
	UserCount int
}

// Mailbox represents an email account
type Mailbox struct {
	Email    string
	Username string
	Domain   string
}

// NewMailService creates a new mail service
func NewMailService(sshClient *SSHClient) *MailService {
	return &MailService{ssh: sshClient}
}

// GetSSHClient returns the underlying SSH client
func (m *MailService) GetSSHClient() *SSHClient {
	return m.ssh
}

// Mail configuration file paths on CMH
const (
	virtualDomainsFile  = "/etc/postfix/virtual_domains"
	virtualMailboxFile  = "/etc/postfix/virtual_mailbox"
	virtualAliasFile    = "/etc/postfix/virtual_alias"
	dovecotUsersFile    = "/etc/dovecot/users"
	virtualMailboxBase  = "/var/mail/vhosts"
)

// ListDomains returns all configured mail domains
func (m *MailService) ListDomains() ([]Domain, error) {
	content, err := m.ssh.ReadFile(virtualDomainsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read domains: %w", err)
	}

	// Get mailbox counts per domain
	mailboxes, err := m.ssh.ReadFile(virtualMailboxFile)
	if err != nil {
		mailboxes = ""
	}

	domainCounts := make(map[string]int)
	for _, line := range strings.Split(mailboxes, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			email := parts[0]
			if at := strings.Index(email, "@"); at > 0 {
				domain := email[at+1:]
				domainCounts[domain]++
			}
		}
	}

	var domains []Domain
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains = append(domains, Domain{
			Name:      line,
			UserCount: domainCounts[line],
		})
	}

	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Name < domains[j].Name
	})

	return domains, nil
}

// AddDomain adds a new mail domain
func (m *MailService) AddDomain(domain string) error {
	// Validate domain format
	if !isValidDomain(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	// Check if domain already exists
	domains, err := m.ListDomains()
	if err != nil {
		return err
	}
	for _, d := range domains {
		if d.Name == domain {
			return fmt.Errorf("domain already exists: %s", domain)
		}
	}

	// Add to virtual_domains
	if err := m.ssh.AppendToFile(virtualDomainsFile, domain); err != nil {
		return fmt.Errorf("failed to add domain: %w", err)
	}

	// Create maildir base for domain
	maildir := fmt.Sprintf("%s/%s", virtualMailboxBase, domain)
	if _, err := m.ssh.Execute(fmt.Sprintf("doas mkdir -p %s && doas chown 5000:5000 %s", maildir, maildir)); err != nil {
		return fmt.Errorf("failed to create maildir: %w", err)
	}

	// Reload postfix
	if _, err := m.ssh.Execute("doas postfix reload"); err != nil {
		return fmt.Errorf("failed to reload postfix: %w", err)
	}

	return nil
}

// DeleteDomain removes a mail domain and all its users
func (m *MailService) DeleteDomain(domain string) error {
	// Get all users for this domain first
	users, err := m.ListMailboxes(domain)
	if err != nil {
		return err
	}

	// Delete all users first
	for _, user := range users {
		if err := m.DeleteMailbox(domain, user.Username); err != nil {
			return fmt.Errorf("failed to delete user %s: %w", user.Email, err)
		}
	}

	// Remove domain from virtual_domains
	if err := m.ssh.DeleteLine(virtualDomainsFile, domain); err != nil {
		return fmt.Errorf("failed to remove domain: %w", err)
	}

	// Remove maildir
	maildir := fmt.Sprintf("%s/%s", virtualMailboxBase, domain)
	if _, err := m.ssh.Execute(fmt.Sprintf("doas rm -rf %s", maildir)); err != nil {
		return fmt.Errorf("failed to remove maildir: %w", err)
	}

	// Reload postfix
	if _, err := m.ssh.Execute("doas postfix reload"); err != nil {
		return fmt.Errorf("failed to reload postfix: %w", err)
	}

	return nil
}

// ListMailboxes returns all mailboxes for a domain
func (m *MailService) ListMailboxes(domain string) ([]Mailbox, error) {
	content, err := m.ssh.ReadFile(virtualMailboxFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read mailboxes: %w", err)
	}

	var mailboxes []Mailbox
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 1 {
			email := parts[0]
			if at := strings.Index(email, "@"); at > 0 {
				mailDomain := email[at+1:]
				if mailDomain == domain {
					mailboxes = append(mailboxes, Mailbox{
						Email:    email,
						Username: email[:at],
						Domain:   mailDomain,
					})
				}
			}
		}
	}

	sort.Slice(mailboxes, func(i, j int) bool {
		return mailboxes[i].Username < mailboxes[j].Username
	})

	return mailboxes, nil
}

// AddMailbox creates a new email account
func (m *MailService) AddMailbox(domain, username, password string) error {
	email := fmt.Sprintf("%s@%s", username, domain)

	// Validate
	if !isValidUsername(username) {
		return fmt.Errorf("invalid username: %s", username)
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Check if user already exists
	existing, err := m.ListMailboxes(domain)
	if err != nil {
		return err
	}
	for _, mb := range existing {
		if mb.Username == username {
			return fmt.Errorf("user already exists: %s", email)
		}
	}

	// Add to postfix virtual_mailbox
	mailboxEntry := fmt.Sprintf("%s    %s/%s/", email, domain, username)
	if err := m.ssh.AppendToFile(virtualMailboxFile, mailboxEntry); err != nil {
		return fmt.Errorf("failed to add mailbox to postfix: %w", err)
	}

	// Regenerate postfix map
	if _, err := m.ssh.Execute(fmt.Sprintf("doas postmap %s", virtualMailboxFile)); err != nil {
		return fmt.Errorf("failed to postmap: %w", err)
	}

	// Create maildir
	maildir := fmt.Sprintf("%s/%s/%s", virtualMailboxBase, domain, username)
	if _, err := m.ssh.Execute(fmt.Sprintf("doas mkdir -p %s && doas chown -R 5000:5000 %s", maildir, maildir)); err != nil {
		return fmt.Errorf("failed to create maildir: %w", err)
	}

	// Add to dovecot users
	dovecotEntry := fmt.Sprintf("%s:{PLAIN}%s", email, password)
	if err := m.ssh.AppendToFile(dovecotUsersFile, dovecotEntry); err != nil {
		return fmt.Errorf("failed to add user to dovecot: %w", err)
	}

	// Reload services
	if _, err := m.ssh.Execute("doas postfix reload && doas doveadm reload"); err != nil {
		return fmt.Errorf("failed to reload services: %w", err)
	}

	return nil
}

// DeleteMailbox removes an email account
func (m *MailService) DeleteMailbox(domain, username string) error {
	email := fmt.Sprintf("%s@%s", username, domain)

	// Remove from dovecot users
	if err := m.ssh.DeleteLine(dovecotUsersFile, fmt.Sprintf("%s:", email)); err != nil {
		return fmt.Errorf("failed to remove from dovecot: %w", err)
	}

	// Remove from postfix virtual_mailbox
	if err := m.ssh.DeleteLine(virtualMailboxFile, email); err != nil {
		return fmt.Errorf("failed to remove from postfix: %w", err)
	}

	// Remove any aliases for this user
	if err := m.ssh.DeleteLine(virtualAliasFile, fmt.Sprintf(".*%s$", email)); err != nil {
		// Ignore alias errors - user might not have any
	}

	// Regenerate postfix maps
	if _, err := m.ssh.Execute(fmt.Sprintf("doas postmap %s", virtualMailboxFile)); err != nil {
		return fmt.Errorf("failed to postmap: %w", err)
	}

	// Optionally remove maildir (commented out to preserve mail)
	// maildir := fmt.Sprintf("%s/%s/%s", virtualMailboxBase, domain, username)
	// m.ssh.Execute(fmt.Sprintf("doas rm -rf %s", maildir))

	// Reload services
	if _, err := m.ssh.Execute("doas postfix reload && doas doveadm reload"); err != nil {
		return fmt.Errorf("failed to reload services: %w", err)
	}

	return nil
}

// ChangePassword updates a user's password
func (m *MailService) ChangePassword(domain, username, newPassword string) error {
	email := fmt.Sprintf("%s@%s", username, domain)

	if len(newPassword) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Update dovecot users file using sed
	newEntry := fmt.Sprintf("%s:{PLAIN}%s", email, newPassword)
	escaped := strings.ReplaceAll(newEntry, "/", "\\/")
	cmd := fmt.Sprintf("doas sed -i 's/^%s:.*/%s/' %s", 
		strings.ReplaceAll(email, "@", "\\@"),
		escaped,
		dovecotUsersFile)
	
	if _, err := m.ssh.Execute(cmd); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Reload dovecot
	if _, err := m.ssh.Execute("doas doveadm reload"); err != nil {
		return fmt.Errorf("failed to reload dovecot: %w", err)
	}

	return nil
}

// TestConnection verifies SSH connectivity to mail server
func (m *MailService) TestConnection() error {
	output, err := m.ssh.Execute("hostname")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	if output == "" {
		return fmt.Errorf("empty response from server")
	}
	return nil
}

// Helper functions

func isValidDomain(domain string) bool {
	if len(domain) < 3 || len(domain) > 253 {
		return false
	}
	if !strings.Contains(domain, ".") {
		return false
	}
	// Basic validation - no spaces, starts with letter/number
	for _, c := range domain {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || 
			(c >= '0' && c <= '9') || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

func isValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 64 {
		return false
	}
	// Allow letters, numbers, dots, underscores, hyphens
	for _, c := range username {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || 
			(c >= '0' && c <= '9') || c == '.' || c == '_' || c == '-') {
			return false
		}
	}
	return true
}
