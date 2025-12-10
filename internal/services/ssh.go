package services

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// SSHClient manages SSH connections to the mail server
type SSHClient struct {
	host        string
	port        int
	user        string
	keyPath     string
	jumpHost    string
	jumpUser    string
	jumpKeyPath string

	mu     sync.Mutex
	client *ssh.Client
}

// SSHConfig holds SSH connection configuration
type SSHConfig struct {
	Host        string
	Port        int
	User        string
	KeyPath     string
	JumpHost    string
	JumpUser    string
	JumpKeyPath string
}

// NewSSHClient creates a new SSH client
func NewSSHClient(cfg SSHConfig) *SSHClient {
	return &SSHClient{
		host:        cfg.Host,
		port:        cfg.Port,
		user:        cfg.User,
		keyPath:     cfg.KeyPath,
		jumpHost:    cfg.JumpHost,
		jumpUser:    cfg.JumpUser,
		jumpKeyPath: cfg.JumpKeyPath,
	}
}

// connect establishes SSH connection (with jump host if configured)
func (c *SSHClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		// Test if connection is still alive
		_, _, err := c.client.SendRequest("keepalive", true, nil)
		if err == nil {
			return nil
		}
		c.client.Close()
		c.client = nil
	}

	// Read private key for target host
	key, err := os.ReadFile(c.keyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect through jump host if configured
	if c.jumpHost != "" {
		// Read jump host key (use separate key if provided, otherwise same key)
		jumpKeyPath := c.jumpKeyPath
		if jumpKeyPath == "" {
			jumpKeyPath = c.keyPath
		}
		
		jumpKey, err := os.ReadFile(jumpKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read jump SSH key: %w", err)
		}

		jumpSigner, err := ssh.ParsePrivateKey(jumpKey)
		if err != nil {
			return fmt.Errorf("failed to parse jump SSH key: %w", err)
		}

		jumpConfig := &ssh.ClientConfig{
			User: c.jumpUser,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(jumpSigner),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		jumpAddr := fmt.Sprintf("%s:22", c.jumpHost)
		jumpClient, err := ssh.Dial("tcp", jumpAddr, jumpConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to jump host: %w", err)
		}

		// Connect to target through jump host
		targetAddr := fmt.Sprintf("%s:%d", c.host, c.port)
		conn, err := jumpClient.Dial("tcp", targetAddr)
		if err != nil {
			jumpClient.Close()
			return fmt.Errorf("failed to dial target through jump: %w", err)
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, targetAddr, config)
		if err != nil {
			conn.Close()
			jumpClient.Close()
			return fmt.Errorf("failed to create client connection: %w", err)
		}

		c.client = ssh.NewClient(ncc, chans, reqs)
	} else {
		// Direct connection
		addr := fmt.Sprintf("%s:%d", c.host, c.port)
		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		c.client = client
	}

	return nil
}

// Execute runs a command on the remote host
func (c *SSHClient) Execute(cmd string) (string, error) {
	if err := c.connect(); err != nil {
		return "", err
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("command failed: %w: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ReadFile reads a file from the remote host
func (c *SSHClient) ReadFile(path string) (string, error) {
	return c.Execute(fmt.Sprintf("doas cat %s 2>/dev/null || cat %s", path, path))
}

// AppendToFile appends content to a file on the remote host
func (c *SSHClient) AppendToFile(path, content string) error {
	// Escape single quotes in content
	escaped := strings.ReplaceAll(content, "'", "'\"'\"'")
	cmd := fmt.Sprintf("echo '%s' | doas tee -a %s > /dev/null", escaped, path)
	_, err := c.Execute(cmd)
	return err
}

// WriteFile writes content to a file (overwrites)
func (c *SSHClient) WriteFile(path, content string) error {
	escaped := strings.ReplaceAll(content, "'", "'\"'\"'")
	cmd := fmt.Sprintf("echo '%s' | doas tee %s > /dev/null", escaped, path)
	_, err := c.Execute(cmd)
	return err
}

// DeleteLine removes a line matching pattern from a file
func (c *SSHClient) DeleteLine(path, pattern string) error {
	// Escape for sed
	escaped := strings.ReplaceAll(pattern, "/", "\\/")
	cmd := fmt.Sprintf("doas sed -i '/^%s/d' %s", escaped, path)
	_, err := c.Execute(cmd)
	return err
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}
