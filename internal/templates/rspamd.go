package templates

// RspamdDashboardHTML contains the Rspamd dashboard template
const RspamdDashboardHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Rspamd Dashboard - MailHub Admin</title>
    <style>
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
            max-width: 1200px;
            margin: 0 auto;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        .card {
            background: white;
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.07);
        }
        .card h2 {
            color: #333;
            font-size: 1.1rem;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
        }
        .card h2 .icon {
            margin-right: 10px;
            font-size: 1.5rem;
        }
        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85rem;
            font-weight: 600;
            margin-left: auto;
        }
        .status-badge.running {
            background: #e8f5e9;
            color: #2e7d32;
        }
        .status-badge.stopped {
            background: #ffebee;
            color: #c62828;
        }
        .metric {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 0;
            border-bottom: 1px solid #eee;
        }
        .metric:last-child {
            border-bottom: none;
        }
        .metric-label {
            color: #666;
            font-size: 0.95rem;
        }
        .metric-value {
            color: #1a73e8;
            font-weight: 600;
            font-size: 1.1rem;
        }
        .chart-container {
            height: 200px;
            margin-bottom: 15px;
        }
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 15px;
        }
        button {
            padding: 8px 16px;
            border: none;
            border-radius: 6px;
            font-size: 0.9rem;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .btn-primary {
            background: #1a73e8;
            color: white;
        }
        .btn-primary:hover {
            background: #1557b0;
        }
        .btn-secondary {
            background: #f0f0f0;
            color: #333;
        }
        .btn-secondary:hover {
            background: #e0e0e0;
        }
        .btn-danger {
            background: #ea4335;
            color: white;
        }
        .btn-danger:hover {
            background: #d33425;
        }
        .whitelist-list {
            list-style: none;
        }
        .whitelist-item {
            padding: 8px;
            background: #f9f9f9;
            border-radius: 4px;
            margin-bottom: 8px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .whitelist-item button {
            padding: 4px 8px;
            font-size: 0.85rem;
        }
        .log-container {
            background: #f5f5f5;
            border-radius: 6px;
            padding: 15px;
            font-family: 'Monaco', 'Courier New', monospace;
            font-size: 0.85rem;
            max-height: 200px;
            overflow-y: auto;
        }
        .log-line {
            padding: 3px 0;
            color: #333;
            border-bottom: 1px solid #eee;
        }
        .log-line.error {
            color: #ea4335;
        }
        .log-line.warn {
            color: #fbbc04;
        }
        .input-group {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
        }
        input[type="text"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 0.95rem;
        }
        .full-width {
            grid-column: 1 / -1;
        }
        .alert {
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        .alert-success {
            background: #e8f5e9;
            color: #2e7d32;
            border-left: 4px solid #2e7d32;
        }
        .alert-error {
            background: #ffebee;
            color: #c62828;
            border-left: 4px solid #c62828;
        }
        .nav-buttons {
            display: flex;
            gap: 10px;
            margin-left: auto;
        }
        .header {
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .header-title h1 {
            color: #1a73e8;
            font-size: 2rem;
            margin-bottom: 5px;
        }
        .header-title p {
            color: #666;
            font-size: 0.95rem;
        }
        .btn-nav {
            padding: 10px 20px;
            background: #1a73e8;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            font-size: 0.9rem;
            cursor: pointer;
            border: none;
            transition: all 0.3s ease;
        }
        .btn-nav:hover {
            background: #1557b0;
        }
        .btn-nav.secondary {
            background: #f0f0f0;
            color: #333;
        }
        .btn-nav.secondary:hover {
            background: #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-title">
                <h1>⚙️ Rspamd Dashboard</h1>
                <p>Antispam and Antivirus Management</p>
            </div>
            <div class="nav-buttons">
                <a href="/" class="btn-nav secondary">← Back to Home</a>
            </div>
        </div>

        <!-- Status Card -->
        <div class="grid">
            <div class="card">
                <h2>
                    <span class="icon">📊</span>
                    Service Status
                    <span class="status-badge running" id="statusBadge">Checking...</span>
                </h2>
                <div id="statusMetrics">
                    <div class="metric">
                        <span class="metric-label">Status</span>
                        <span class="metric-value" id="statusValue">Checking...</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Version</span>
                        <span class="metric-value" id="versionValue">--</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Memory Usage</span>
                        <span class="metric-value" id="memoryValue">--</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">CPU Usage</span>
                        <span class="metric-value" id="cpuValue">--</span>
                    </div>
                </div>
                <div class="button-group">
                    <button class="btn-primary" onclick="startService()">Start</button>
                    <button class="btn-secondary" onclick="restartService()">Restart</button>
                    <button class="btn-danger" onclick="stopService()">Stop</button>
                </div>
            </div>

            <!-- Metrics Card -->
            <div class="card">
                <h2>
                    <span class="icon">📈</span>
                    Performance Metrics
                </h2>
                <div id="metricsData">
                    <div class="metric">
                        <span class="metric-label">Messages Scanned</span>
                        <span class="metric-value" id="messagesValue">0</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Spam Detected</span>
                        <span class="metric-value" id="spamValue">0</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Spam Percentage</span>
                        <span class="metric-value" id="spamPercentValue">0%</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Average Score</span>
                        <span class="metric-value" id="scoreValue">--</span>
                    </div>
                </div>
                <button class="btn-primary" style="width: 100%; margin-top: 15px;" onclick="refreshMetrics()">Refresh</button>
            </div>

            <!-- Configuration Card -->
            <div class="card">
                <h2>
                    <span class="icon">⚙️</span>
                    Configuration
                </h2>
                <div id="configData">
                    <div class="metric">
                        <span class="metric-label">Worker Max Tasks</span>
                        <span class="metric-value" id="maxTasksValue">20</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Worker Count</span>
                        <span class="metric-value" id="workerCountValue">1</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Worker Timeout</span>
                        <span class="metric-value" id="timeoutValue">30s</span>
                    </div>
                    <div class="metric">
                        <span class="metric-label">Redis Memory</span>
                        <span class="metric-value" id="redisMemValue">256mb</span>
                    </div>
                </div>
                <button class="btn-primary" style="width: 100%; margin-top: 15px;" onclick="openConfigEditor()">Edit</button>
            </div>
        </div>

        <!-- Whitelist Card -->
        <div class="card full-width">
            <h2>
                <span class="icon">✅</span>
                SPF Whitelist Management
            </h2>
            <div class="input-group">
                <input type="text" id="whitelistInput" placeholder="Enter domain, email, or IP to whitelist (e.g., *.example.com, user@example.com, 192.168.1.1)">
                <button class="btn-primary" onclick="addToWhitelist()">Add</button>
            </div>
            <div id="whitelistContainer">
                <p style="color: #999; text-align: center;">Loading whitelist...</p>
            </div>
        </div>

        <!-- Logs Card -->
        <div class="card full-width">
            <h2>
                <span class="icon">📝</span>
                Recent Activity Logs
            </h2>
            <div id="logsContainer" class="log-container">
                <div class="log-line">Loading logs...</div>
            </div>
            <button class="btn-secondary" style="width: 100%; margin-top: 15px;" onclick="refreshLogs()">Refresh Logs</button>
        </div>

        <!-- Export Card -->
        <div class="card full-width">
            <h2>
                <span class="icon">📥</span>
                Data Export & Advanced
            </h2>
            <div class="button-group">
                <button class="btn-primary" onclick="exportMetrics()">Export as JSON</button>
                <button class="btn-secondary" onclick="viewDocs()">Documentation</button>
            </div>
        </div>
    </div>

    <script>
        const API_BASE = '/rspamd';

        async function fetchStatus() {
            try {
                const response = await fetch(API_BASE + '/status');
                const data = await response.json();
                if (data.success) {
                    const status = data.data;
                    document.getElementById('statusValue').textContent = status.is_running ? 'Running' : 'Stopped';
                    document.getElementById('statusBadge').className = 'status-badge ' + (status.is_running ? 'running' : 'stopped');
                    document.getElementById('versionValue').textContent = status.version || 'N/A';
                    document.getElementById('memoryValue').textContent = status.memory || '--';
                    document.getElementById('cpuValue').textContent = status.cpu || '--';
                }
            } catch (error) {
                console.error('Error fetching status:', error);
            }
        }

        async function fetchMetrics() {
            try {
                const response = await fetch(API_BASE + '/metrics');
                const data = await response.json();
                if (data.success) {
                    const metrics = data.data;
                    document.getElementById('messagesValue').textContent = metrics.message_count || '0';
                    document.getElementById('spamValue').textContent = metrics.spam_count || '0';
                    document.getElementById('spamPercentValue').textContent = (metrics.spam_percentage || 0).toFixed(2) + '%';
                    document.getElementById('scoreValue').textContent = (metrics.average_score || 0).toFixed(2);
                }
            } catch (error) {
                console.error('Error fetching metrics:', error);
            }
        }

        async function fetchConfig() {
            try {
                const response = await fetch(API_BASE + '/config');
                const data = await response.json();
                if (data.success) {
                    const config = data.data;
                    document.getElementById('maxTasksValue').textContent = config.worker_max_tasks;
                    document.getElementById('workerCountValue').textContent = config.worker_count;
                    document.getElementById('timeoutValue').textContent = config.worker_timeout + 's';
                    document.getElementById('redisMemValue').textContent = config.redis_memory;
                }
            } catch (error) {
                console.error('Error fetching config:', error);
            }
        }

        async function fetchWhitelist() {
            try {
                const response = await fetch(API_BASE + '/whitelist');
                const data = await response.json();
                if (data.success) {
                    const whitelist = data.data;
                    const container = document.getElementById('whitelistContainer');
                    let html = '<ul class="whitelist-list">';
                    
                    const items = [
                        ...(whitelist.domains || []),
                        ...(whitelist.emails || []),
                        ...(whitelist.ips || [])
                    ];
                    
                    if (items.length === 0) {
                        html += '<li style="text-align: center; color: #999; padding: 15px;">No entries in whitelist</li>';
                    } else {
                        items.forEach(item => {
                            html += '<li class="whitelist-item"><span>' + item + '</span><button class="btn-danger" onclick="removeFromWhitelist(\'' + item.replace(/'/g, "\\'") + '\')">Remove</button></li>';
                        });
                    }
                    html += '</ul>';
                    container.innerHTML = html;
                }
            } catch (error) {
                console.error('Error fetching whitelist:', error);
            }
        }

        async function fetchLogs() {
            try {
                const response = await fetch(API_BASE + '/logs?lines=20');
                const data = await response.json();
                if (data.success) {
                    const logs = data.data;
                    const container = document.getElementById('logsContainer');
                    let html = '';
                    logs.forEach(log => {
                        const className = log.level === 'ERROR' ? 'error' : (log.level === 'WARN' ? 'warn' : '');
                        html += '<div class="log-line ' + className + '">' + log.message + '</div>';
                    });
                    container.innerHTML = html || '<div class="log-line">No logs available</div>';
                }
            } catch (error) {
                console.error('Error fetching logs:', error);
            }
        }

        async function startService() {
            if (confirm('Start Rspamd service?')) {
                try {
                    const response = await fetch(API_BASE + '/service/start', { method: 'POST' });
                    const data = await response.json();
                    alert(data.message || (data.success ? 'Service started' : 'Error'));
                    fetchStatus();
                } catch (error) {
                    alert('Error: ' + error.message);
                }
            }
        }

        async function stopService() {
            if (confirm('Stop Rspamd service? This will stop spam filtering.')) {
                try {
                    const response = await fetch(API_BASE + '/service/stop', { method: 'POST' });
                    const data = await response.json();
                    alert(data.message || (data.success ? 'Service stopped' : 'Error'));
                    fetchStatus();
                } catch (error) {
                    alert('Error: ' + error.message);
                }
            }
        }

        async function restartService() {
            if (confirm('Restart Rspamd service?')) {
                try {
                    const response = await fetch(API_BASE + '/service/restart', { method: 'POST' });
                    const data = await response.json();
                    alert(data.message || (data.success ? 'Service restarted' : 'Error'));
                    fetchStatus();
                } catch (error) {
                    alert('Error: ' + error.message);
                }
            }
        }

        async function addToWhitelist() {
            const entry = document.getElementById('whitelistInput').value.trim();
            if (!entry) {
                alert('Please enter a domain, email, or IP');
                return;
            }
            try {
                const response = await fetch(API_BASE + '/whitelist', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ entry: entry })
                });
                const data = await response.json();
                if (data.success) {
                    document.getElementById('whitelistInput').value = '';
                    fetchWhitelist();
                    alert('Entry added to whitelist');
                } else {
                    alert('Error: ' + data.error);
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        }

        async function removeFromWhitelist(entry) {
            if (confirm('Remove "' + entry + '" from whitelist?')) {
                try {
                    const response = await fetch(API_BASE + '/whitelist', {
                        method: 'DELETE',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ entry: entry })
                    });
                    const data = await response.json();
                    if (data.success) {
                        fetchWhitelist();
                    } else {
                        alert('Error: ' + data.error);
                    }
                } catch (error) {
                    alert('Error: ' + error.message);
                }
            }
        }

        async function exportMetrics() {
            try {
                const response = await fetch(API_BASE + '/export');
                const data = await response.json();
                const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'rspamd-metrics-' + new Date().toISOString().split('T')[0] + '.json';
                a.click();
            } catch (error) {
                alert('Error: ' + error.message);
            }
        }

        function openConfigEditor() {
            alert('Configuration editor coming soon');
        }

        function viewDocs() {
            window.open('https://rspamd.com/doc/', '_blank');
        }

        function refreshMetrics() {
            fetchMetrics();
            fetchLogs();
        }

        function refreshLogs() {
            fetchLogs();
        }

        // Initial load
        document.addEventListener('DOMContentLoaded', function() {
            fetchStatus();
            fetchMetrics();
            fetchConfig();
            fetchWhitelist();
            fetchLogs();

            // Refresh every 30 seconds
            setInterval(() => {
                fetchStatus();
                fetchMetrics();
                fetchLogs();
            }, 30000);
        });
    </script>
</body>
</html>
`
