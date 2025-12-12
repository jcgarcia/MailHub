package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ingasti/mailhub-admin/internal/services"
	"github.com/Ingasti/mailhub-admin/internal/templates"
)

// HandleRspamdDashboard renders the Rspamd dashboard page
func HandleRspamdDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(templates.RspamdDashboardHTML))
}

// RspamdResponse wraps API responses
type RspamdResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// HandleRspamdStatus returns the current Rspamd status
func HandleRspamdStatus(w http.ResponseWriter, r *http.Request) {
	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	status, err := rspamd.GetStatus()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Data:    status,
	})
}

// HandleRspamdMetrics returns Rspamd performance metrics
func HandleRspamdMetrics(w http.ResponseWriter, r *http.Request) {
	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	metrics, err := rspamd.GetMetrics()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Data:    metrics,
	})
}

// HandleRspamdConfig returns the current Rspamd configuration
func HandleRspamdConfig(w http.ResponseWriter, r *http.Request) {
	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	config, err := rspamd.GetConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Data:    config,
	})
}

// HandleRspamdConfigUpdate updates Rspamd configuration
func HandleRspamdConfigUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var config services.RspamdConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.UpdateConfig(&config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Configuration updated successfully",
	})
}

// HandleRspamdWhitelist returns the SPF whitelist
func HandleRspamdWhitelist(w http.ResponseWriter, r *http.Request) {
	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	whitelist, err := rspamd.GetWhitelist()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Data:    whitelist,
	})
}

// HandleRspamdWhitelistAdd adds an entry to the whitelist
func HandleRspamdWhitelistAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var req struct {
		Entry string `json:"entry"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.AddToWhitelist(req.Entry); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Entry added to whitelist",
	})
}

// HandleRspamdWhitelistRemove removes an entry from the whitelist
func HandleRspamdWhitelistRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var req struct {
		Entry string `json:"entry"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.RemoveFromWhitelist(req.Entry); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Entry removed from whitelist",
	})
}

// HandleRspamdLogs returns recent Rspamd logs
func HandleRspamdLogs(w http.ResponseWriter, r *http.Request) {
	// Get line count from query params
	lineStr := r.URL.Query().Get("lines")
	lines := 50
	if lineStr != "" {
		if l, err := strconv.Atoi(lineStr); err == nil && l > 0 {
			lines = l
		}
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	logs, err := rspamd.GetLogs(lines)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Data:    logs,
	})
}

// HandleRspamdServiceStart starts the Rspamd service
func HandleRspamdServiceStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.StartService(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Rspamd service started",
	})
}

// HandleRspamdServiceStop stops the Rspamd service
func HandleRspamdServiceStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.StopService(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Rspamd service stopped",
	})
}

// HandleRspamdServiceRestart restarts the Rspamd service
func HandleRspamdServiceRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	if err := rspamd.RestartService(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RspamdResponse{
		Success: true,
		Message: "Rspamd service restarted",
	})
}

// HandleRspamdExport exports all metrics as JSON
func HandleRspamdExport(w http.ResponseWriter, r *http.Request) {
	h := GetHandler()
	rspamd := services.NewRspamdService(h.Mail.GetSSHClient())

	data, err := rspamd.ExportMetricsJSON()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RspamdResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
