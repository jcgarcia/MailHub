package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RspamdService provides Rspamd management operations
type RspamdService struct {
	ssh *SSHClient
}

// RspamdStatus represents the current status of Rspamd
type RspamdStatus struct {
	IsRunning   bool   `json:"is_running"`
	Version     string `json:"version"`
	Memory      string `json:"memory"`
	CPU         string `json:"cpu"`
	ProcessID   int    `json:"process_id"`
	Uptime      string `json:"uptime"`
	MessageRate string `json:"message_rate"`
}

// RspamdMetrics represents Rspamd performance metrics
type RspamdMetrics struct {
	MessageCount      int64   `json:"message_count"`
	SpamCount         int64   `json:"spam_count"`
	HamCount          int64   `json:"ham_count"`
	SpamPercentage    float64 `json:"spam_percentage"`
	AverageScore      float64 `json:"average_score"`
	LearnedSpam       int64   `json:"learned_spam"`
	LearnedHam        int64   `json:"learned_ham"`
	FuzzyMatches      int64   `json:"fuzzy_matches"`
	DNSBLMatches      int64   `json:"dnsbl_matches"`
}

// RspamdConfig represents Rspamd configuration
type RspamdConfig struct {
	WorkerMaxTasks  int    `json:"worker_max_tasks"`
	WorkerCount     int    `json:"worker_count"`
	WorkerTimeout   int    `json:"worker_timeout"`
	RedisMemory     string `json:"redis_memory"`
	SPFEnabled      bool   `json:"spf_enabled"`
	DKIMEnabled     bool   `json:"dkim_enabled"`
	SURBLEnabled    bool   `json:"surbl_enabled"`
	FuzzyEnabled    bool   `json:"fuzzy_enabled"`
	CharTableEnable bool   `json:"chart_enable"`
}

// RspamdWhitelist represents whitelisted senders
type RspamdWhitelist struct {
	Domains []string `json:"domains"`
	IPs     []string `json:"ips"`
	Emails  []string `json:"emails"`
}

// RspamdLog represents a log entry
type RspamdLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Score     string `json:"score,omitempty"`
}

// Config file paths on CMH
const (
	rspamdWorkerConf   = "/etc/rspamd/local.d/worker-normal.conf"
	rspamdOptionsConf  = "/etc/rspamd/local.d/options.inc"
	rspamdModulesConf  = "/etc/rspamd/local.d/modules.conf"
	rspamdWhitelistTxt = "/etc/rspamd/spf_whitelist.txt"
	redisConf          = "/etc/redis.conf"
	rspamdLogFile      = "/var/log/rspamd/rspamd.log"
)

// NewRspamdService creates a new Rspamd service
func NewRspamdService(sshClient *SSHClient) *RspamdService {
	return &RspamdService{ssh: sshClient}
}

// GetStatus returns the current status of Rspamd
func (r *RspamdService) GetStatus() (*RspamdStatus, error) {
	// Check if Rspamd is running
	cmd := "doas rc-service rspamd status && echo 'RUNNING' || echo 'STOPPED'"
	output, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to check Rspamd status: %w", err)
	}

	status := &RspamdStatus{
		IsRunning: strings.Contains(output, "started"),
	}

	if !status.IsRunning {
		return status, nil
	}

	// Get version
	versionCmd := "rspamd --version | head -1"
	versionOut, _ := r.ssh.ExecuteCommand(versionCmd)
	status.Version = strings.TrimSpace(versionOut)

	// Get process info
	psCmd := "doas ps aux | grep '[r]spawmd' | head -1 | awk '{print $2, $6}'"
	psOut, _ := r.ssh.ExecuteCommand(psCmd)
	if psOut != "" {
		parts := strings.Fields(psOut)
		if len(parts) >= 2 {
			fmt.Sscanf(parts[0], "%d", &status.ProcessID)
			status.Memory = parts[1]
		}
	}

	// Get CPU usage from top
	topCmd := "doas top -bn 1 | grep -E '^[%]|rspamd' | tail -1 | awk '{print $9}'"
	topOut, _ := r.ssh.ExecuteCommand(topCmd)
	status.CPU = strings.TrimSpace(topOut) + "%"

	return status, nil
}

// GetMetrics returns Rspamd metrics
func (r *RspamdService) GetMetrics() (*RspamdMetrics, error) {
	// Try to get metrics from Rspamd HTTP interface
	cmd := `doas wget -q -O - http://127.0.0.1:11334/stat | grep -E '"(scanned|spam|ham|score)"|Total:' | head -20`
	output, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		// Fallback: parse from logs
		return r.getMetricsFromLogs()
	}

	metrics := &RspamdMetrics{}

	// Parse output (simplified parsing for common metrics)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "scanned") {
			fmt.Sscanf(line, "%d", &metrics.MessageCount)
		}
		if strings.Contains(line, "spam") && !strings.Contains(line, "learned") {
			fmt.Sscanf(line, "%d", &metrics.SpamCount)
		}
		if strings.Contains(line, "ham") {
			fmt.Sscanf(line, "%d", &metrics.HamCount)
		}
	}

	// Calculate spam percentage
	if metrics.MessageCount > 0 {
		metrics.SpamPercentage = float64(metrics.SpamCount) / float64(metrics.MessageCount) * 100
	}

	return metrics, nil
}

// getMetricsFromLogs extracts metrics from log file
func (r *RspamdService) getMetricsFromLogs() (*RspamdMetrics, error) {
	metrics := &RspamdMetrics{}

	// Get last 1000 log lines
	cmd := fmt.Sprintf("doas tail -1000 %s", rspamdLogFile)
	output, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return metrics, fmt.Errorf("failed to read Rspamd logs: %w", err)
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "spam") {
			metrics.SpamCount++
		}
		if strings.Contains(line, "ham") {
			metrics.HamCount++
		}
		metrics.MessageCount++
	}

	if metrics.MessageCount > 0 {
		metrics.SpamPercentage = float64(metrics.SpamCount) / float64(metrics.MessageCount) * 100
	}

	return metrics, nil
}

// GetConfig returns current Rspamd configuration
func (r *RspamdService) GetConfig() (*RspamdConfig, error) {
	config := &RspamdConfig{
		WorkerMaxTasks:  20,
		WorkerCount:     1,
		WorkerTimeout:   30,
		RedisMemory:     "256mb",
		SPFEnabled:      true,
		DKIMEnabled:     true,
		SURBLEnabled:    true,
		FuzzyEnabled:    true,
		CharTableEnable: true,
	}

	// Read worker config
	workerContent, err := r.ssh.ReadFile(rspamdWorkerConf)
	if err == nil {
		r.parseWorkerConfig(workerContent, config)
	}

	// Read options config
	optionsContent, err := r.ssh.ReadFile(rspamdOptionsConf)
	if err == nil {
		r.parseOptionsConfig(optionsContent, config)
	}

	// Read Redis config
	redisContent, err := r.ssh.ReadFile(redisConf)
	if err == nil {
		r.parseRedisConfig(redisContent, config)
	}

	return config, nil
}

// parseWorkerConfig extracts worker settings from config file
func (r *RspamdService) parseWorkerConfig(content string, config *RspamdConfig) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "max_tasks") {
			fmt.Sscanf(line, "max_tasks = %d", &config.WorkerMaxTasks)
		}
		if strings.HasPrefix(line, "count") {
			fmt.Sscanf(line, "count = %d", &config.WorkerCount)
		}
		if strings.HasPrefix(line, "timeout") {
			fmt.Sscanf(line, "timeout = %d", &config.WorkerTimeout)
		}
	}
}

// parseOptionsConfig extracts options from config file
func (r *RspamdService) parseOptionsConfig(content string, config *RspamdConfig) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "disable_spf") && strings.Contains(line, "yes") {
			config.SPFEnabled = false
		}
		if strings.Contains(line, "disable_dkim") && strings.Contains(line, "yes") {
			config.DKIMEnabled = false
		}
	}
}

// parseRedisConfig extracts Redis memory limits
func (r *RspamdService) parseRedisConfig(content string, config *RspamdConfig) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "maxmemory ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				config.RedisMemory = parts[1]
			}
		}
	}
}

// UpdateConfig updates Rspamd configuration
func (r *RspamdService) UpdateConfig(config *RspamdConfig) error {
	// Update worker config
	workerConf := fmt.Sprintf(`
# Worker configuration - Optimized for Celeron CPU
worker "normal" {
  max_tasks = %d;
  count = %d;
  timeout = %ds;
}
`, config.WorkerMaxTasks, config.WorkerCount, config.WorkerTimeout)

	if err := r.ssh.WriteFile(rspamdWorkerConf, workerConf); err != nil {
		return fmt.Errorf("failed to update worker config: %w", err)
	}

	// Reload Rspamd
	cmd := "doas rc-service rspamd reload"
	if _, err := r.ssh.ExecuteCommand(cmd); err != nil {
		return fmt.Errorf("failed to reload Rspamd: %w", err)
	}

	return nil
}

// GetWhitelist returns the current SPF whitelist
func (r *RspamdService) GetWhitelist() (*RspamdWhitelist, error) {
	content, err := r.ssh.ReadFile(rspamdWhitelistTxt)
	if err != nil {
		return nil, fmt.Errorf("failed to read whitelist: %w", err)
	}

	whitelist := &RspamdWhitelist{
		Domains: []string{},
		IPs:     []string{},
		Emails:  []string{},
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "*") || strings.HasPrefix(line, "*.") {
			whitelist.Domains = append(whitelist.Domains, line)
		} else if strings.Contains(line, "@") {
			whitelist.Emails = append(whitelist.Emails, line)
		} else if isValidIP(line) {
			whitelist.IPs = append(whitelist.IPs, line)
		}
	}

	return whitelist, nil
}

// AddToWhitelist adds a sender to the whitelist
func (r *RspamdService) AddToWhitelist(entry string) error {
	if entry == "" {
		return fmt.Errorf("whitelist entry cannot be empty")
	}

	// Read current whitelist
	whitelist, err := r.GetWhitelist()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, d := range whitelist.Domains {
		if d == entry {
			return fmt.Errorf("entry already in whitelist")
		}
	}
	for _, e := range whitelist.Emails {
		if e == entry {
			return fmt.Errorf("entry already in whitelist")
		}
	}
	for _, ip := range whitelist.IPs {
		if ip == entry {
			return fmt.Errorf("entry already in whitelist")
		}
	}

	// Append to whitelist file
	cmd := fmt.Sprintf("echo '%s' | doas tee -a %s", entry, rspamdWhitelistTxt)
	if _, err := r.ssh.ExecuteCommand(cmd); err != nil {
		return fmt.Errorf("failed to add to whitelist: %w", err)
	}

	return nil
}

// RemoveFromWhitelist removes a sender from the whitelist
func (r *RspamdService) RemoveFromWhitelist(entry string) error {
	if entry == "" {
		return fmt.Errorf("whitelist entry cannot be empty")
	}

	// Use sed to remove the entry
	cmd := fmt.Sprintf("doas sed -i '/%s/d' %s", entry, rspamdWhitelistTxt)
	if _, err := r.ssh.ExecuteCommand(cmd); err != nil {
		return fmt.Errorf("failed to remove from whitelist: %w", err)
	}

	return nil
}

// GetLogs returns recent Rspamd logs
func (r *RspamdService) GetLogs(lines int) ([]RspamdLog, error) {
	if lines <= 0 {
		lines = 50
	}

	cmd := fmt.Sprintf("doas tail -%d %s", lines, rspamdLogFile)
	output, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	var logs []RspamdLog
	logLines := strings.Split(output, "\n")
	for _, line := range logLines {
		if line == "" {
			continue
		}
		log := RspamdLog{
			Message: line,
		}

		// Try to parse timestamp
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			log.Timestamp = parts[0] + " " + parts[1]
		}

		// Determine log level
		if strings.Contains(line, "ERROR") || strings.Contains(line, "error") {
			log.Level = "ERROR"
		} else if strings.Contains(line, "WARN") || strings.Contains(line, "warn") {
			log.Level = "WARN"
		} else {
			log.Level = "INFO"
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// TestConnection tests the Rspamd connection
func (r *RspamdService) TestConnection() error {
	cmd := "doas rc-service rspamd status | grep -q 'started'"
	_, err := r.ssh.ExecuteCommand(cmd)
	return err
}

// isValidIP checks if a string is a valid IP address
func isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 && len(parts) != 2 {
		return false
	}
	// Simple validation
	if len(parts) == 4 {
		for _, part := range parts {
			if part == "" || part == "*" {
				continue
			}
			// Additional validation could go here
		}
	}
	return true
}

// RestartService restarts the Rspamd service
func (r *RspamdService) RestartService() error {
	cmd := "doas rc-service rspamd restart"
	_, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to restart Rspamd: %w", err)
	}
	return nil
}

// StopService stops the Rspamd service
func (r *RspamdService) StopService() error {
	cmd := "doas rc-service rspamd stop"
	_, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to stop Rspamd: %w", err)
	}
	return nil
}

// StartService starts the Rspamd service
func (r *RspamdService) StartService() error {
	cmd := "doas rc-service rspamd start"
	_, err := r.ssh.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to start Rspamd: %w", err)
	}
	return nil
}

// ExportMetricsJSON exports current metrics as JSON
func (r *RspamdService) ExportMetricsJSON() ([]byte, error) {
	status, err := r.GetStatus()
	if err != nil {
		return nil, err
	}

	metrics, err := r.GetMetrics()
	if err != nil {
		return nil, err
	}

	config, err := r.GetConfig()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"status":  status,
		"metrics": metrics,
		"config":  config,
	}

	return json.MarshalIndent(data, "", "  ")
}
