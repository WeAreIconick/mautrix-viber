// Package webadmin provides web admin panel for bridge management and statistics.
package webadmin

import (
	"net/http"

	"github.com/example/mautrix-viber/internal/api"
)

// AdminHandler handles web admin requests.
type AdminHandler struct {
	// Additional fields for admin operations
}

// NewAdminHandler creates a new admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ServeAdminPanel serves the admin panel HTML.
func (h *AdminHandler) ServeAdminPanel(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>mautrix-viber Admin Panel</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		.stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
		.stat-card { border: 1px solid #ddd; padding: 15px; border-radius: 5px; }
		.stat-value { font-size: 2em; font-weight: bold; }
		.stat-label { color: #666; }
	</style>
</head>
<body>
	<h1>mautrix-viber Bridge Admin Panel</h1>
	<div id="stats" class="stats">
		<div class="stat-card">
			<div class="stat-value" id="messages-bridged">0</div>
			<div class="stat-label">Messages Bridged</div>
		</div>
		<div class="stat-card">
			<div class="stat-value" id="users-linked">0</div>
			<div class="stat-label">Users Linked</div>
		</div>
		<div class="stat-card">
			<div class="stat-value" id="rooms-mapped">0</div>
			<div class="stat-label">Rooms Mapped</div>
		</div>
		<div class="stat-card">
			<div class="stat-value" id="uptime">0</div>
			<div class="stat-label">Uptime</div>
		</div>
	</div>
	<script>
		async function updateStats() {
			const response = await fetch('/api/info');
			const data = await response.json();
			document.getElementById('messages-bridged').textContent = data.statistics.messages_bridged;
			document.getElementById('users-linked').textContent = data.statistics.users_linked;
			document.getElementById('rooms-mapped').textContent = data.statistics.rooms_mapped;
			document.getElementById('uptime').textContent = data.uptime;
		}
		updateStats();
		setInterval(updateStats, 5000);
	</script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// HandleAPI handles admin API requests.
func (h *AdminHandler) HandleAPI(w http.ResponseWriter, r *http.Request) {
	// Return bridge info as JSON
	w.Header().Set("Content-Type", "application/json")
	api.InfoHandler(w, r)
}

