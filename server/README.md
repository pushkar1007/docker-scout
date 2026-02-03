Docker Scout Server

This server exposes a REST API + real-time WebSocket streams for Docker resources.

**Default server URL:** `http://localhost:8089`
**WebSocket endpoints:**
- Stats: `ws://localhost:8089/stats`
- Terminal: `ws://localhost:8089/terminal`

## Run (Linux/macOS)
From `server/`:

```bash
# Clone the repository
cd server

# Start the server
docker compose up -d

# Verify it's running
curl http://localhost:8089/system/health
```

**Server is now running at:** `http://localhost:8089`

### Stop the Server

```bash
docker compose down
```

## Run (Windows)
From `server\` (Git Bash):

```bash
docker compose logs -f docker_api_server
```

---

## Configuration

### Environment Variables

```bash
# Set custom port (default: 8089)
export DOCKER_SCOUT_ADDR=:9000
docker compose up -d
```

**Environment File (.env):**
```env
DOCKER_SCOUT_ADDR=:8089
```

## Run via Docker Compose (Linux)
```bash
docker compose up --build
```

## Run via Docker Compose (Windows)
Edit `server/docker-compose.yaml` to mount the Windows named pipe:

```yaml
volumes:
  - \\.\pipe\docker_engine:/var/run/docker.sock
```

Then:

```bash
docker compose up --build
```

## Test Containers (Fixed + Random Usage)

### Fixed usage containers (low steady load)
```bash
docker run -d --rm --name scout-fixed-1 alpine:3.20 sh -c "while true; do sleep 1; done"
docker run -d --rm --name scout-fixed-2 alpine:3.20 sh -c "while true; do sleep 1; done"
docker run -d --rm --name scout-fixed-3 alpine:3.20 sh -c "while true; do sleep 1; done"
```

### Random usage container (memory 0–200MB, CPU ~1–5%)
```bash
docker run -d --rm --name scout-random-load python:3.12-alpine sh -c 'python - << "PY"
import random, time

def burn_cpu(percent, duration=1.0):
    busy = duration * (percent / 100.0)
    idle = duration - busy
    end = time.time() + busy
    while time.time() < end:
        pass
    time.sleep(max(idle, 0))

buf = bytearray()
while True:
    target_mb = random.randint(0, 200)
    buf = bytearray(target_mb * 1024 * 1024)
    burn_cpu(random.randint(1, 5), 1.0)
PY'
```

<<<<<<< HEAD
## Routes And Curl Commands (All Endpoints)

Set a base URL once:

```bash
# Direct
BASE=http://localhost:8089
# Proxy
# BASE=http://localhost/scout-server
```

### UI
```bash
curl "$BASE/"
curl "$BASE/echo"
curl "$BASE/index"
```

### System
```bash
curl "$BASE/system/health"
curl -X POST "$BASE/system/nuke?confirm=true"
```

### Containers
```bash
curl "$BASE/containers"
curl -X POST "$BASE/containers/create?start=true" \
  -H "Content-Type: application/json" \
  -d '{"name":"scout-demo","image":"alpine:3.20","cmd":["sh","-c","sleep 600"]}'
curl -X POST "$BASE/containers/start?id=<container_id>"
curl -X POST "$BASE/containers/stop?id=<container_id>"
curl -X POST "$BASE/containers/pause?id=<container_id>"
curl -X POST "$BASE/containers/unpause?id=<container_id>"
curl -X DELETE "$BASE/containers/remove?id=<container_id>"
```

### Images
```bash
curl "$BASE/images"
curl -X DELETE "$BASE/images/remove?name=nginx:latest&force=true"
curl -X POST "$BASE/images/build" \
  -H "Content-Type: application/json" \
  -d '{"context_path":"/path/to/context","dockerfile":"Dockerfile","tag":"myapp:dev","no_cache":false}'
```

### Volumes
```bash
curl "$BASE/volumes"
curl "$BASE/volumes/inspect?name=<volume_name>"
curl -X POST "$BASE/volumes/create" \
  -H "Content-Type: application/json" \
  -d '{"name":"scout-vol","driver":"local"}'
curl -X DELETE "$BASE/volumes/remove?name=<volume_name>"
curl -X POST "$BASE/volumes/prune"
```

### Networks
```bash
curl "$BASE/networks"
curl -X POST "$BASE/networks/create" \
  -H "Content-Type: application/json" \
  -d '{"name":"scout-net","driver":"bridge","options":{}}'
curl -X DELETE "$BASE/networks/remove?id=<network_id_or_name>"
```

### Stats
```bash
curl "$BASE/stats"
```

### Events (WebSocket + Publish)
```bash
curl -X POST "$BASE/events/publish" \
  -H "Content-Type: application/json" \
  -d '{"event":"custom_event","data":"hello"}'

# WebSocket handshake only (curl can’t send frames)
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" -H "Sec-WebSocket-Version: 13" \
  "$BASE/events"
```

### Terminal (WebSocket)
```bash
# WebSocket handshake only (curl can’t send frames)
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" -H "Sec-WebSocket-Version: 13" \
  "$BASE/terminal"
=======
## Quick Start WebSocket Examples

### Monitor Docker Stats in Real-Time

```html
<!DOCTYPE html>
<html>
<body>
  <div id="stats"></div>
  <script>
    const ws = new WebSocket('ws://localhost:8089/stats');
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      document.getElementById('stats').innerHTML = 
        '<pre>' + JSON.stringify(data.stats, null, 2) + '</pre>';
    };
  </script>
</body>
</html>
```

### Execute Commands via Terminal WebSocket

```html
<!DOCTYPE html>
<html>
<body>
  <input id="cmd" placeholder="Enter command" />
  <button onclick="runCommand()">Run</button>
  <button onclick="stopCommand()">Stop</button>
  <pre id="output"></pre>
  
  <script>
    let ws;
    
    function runCommand() {
      if (ws) ws.close();
      ws = new WebSocket('ws://localhost:8089/terminal');
      
      ws.onopen = () => {
        const cmd = document.getElementById('cmd').value;
        ws.send(JSON.stringify({ type: 'start', command: cmd }));
      };
      
      ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        const output = document.getElementById('output');
        
        if (msg.type === 'stdout' || msg.type === 'stderr') {
          output.textContent += msg.data + '\\n';
        } else if (msg.type === 'exit') {
          output.textContent += `\\nExited with code: ${msg.code}\\n`;
        }
      };
    }
    
    function stopCommand() {
      if (ws) ws.send(JSON.stringify({ type: 'exitTerm' }));
    }
  </script>
</body>
</html>
>>>>>>> 22f7920 (Updated WebSocket)
```

## API Requests (curl)

### System Management
```bash
curl -X POST "http://localhost:8089/system/nuke?confirm=true"
```

### Containers
```bash
curl http://localhost:8089/containers
curl -X POST "http://localhost:8089/containers/start?id=<container_id>"
curl -X POST "http://localhost:8089/containers/stop?id=<container_id>"
curl -X POST "http://localhost:8089/containers/pause?id=<container_id>"
curl -X POST "http://localhost:8089/containers/unpause?id=<container_id>"
curl -X DELETE "http://localhost:8089/containers/remove?id=<container_id>"
```

### Images
```bash
curl http://localhost:8089/images
```

---

### Remove Image

**DELETE** `/images/remove?name=<image_name>&force=<true|false>`

Removes an image.

**Parameters:**
- `name` (required): Image name or tag
- `force` (optional): Force removal (`"true"` or `"false"`)

**Response:**
```json
{
  "status": "removed",
  "image": "nginx:latest"
}
```

**Example:**
```bash
curl -X DELETE "http://localhost:8089/images/remove?name=nginx:latest"
curl -X DELETE "http://localhost:8089/images/remove?name=nginx:latest&force=true"
```

---

### Build Image

**POST** `/images/build`

Builds a new image from Dockerfile. Returns streaming build output.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "context_path": "/path/to/dockerfile/directory",
  "dockerfile": "Dockerfile",
  "tag": "myapp:1.0.0",
  "no_cache": false
}
```

**Parameters:**
- `context_path` (required): Path to directory containing Dockerfile
- `dockerfile` (optional): Dockerfile name (default: "Dockerfile")
- `tag` (required): Image tag
- `no_cache` (optional): Disable layer caching

**Response:**
```
Streaming build output:
Step 1/5 : FROM golang:1.21
 ---> abc123def456
Step 2/5 : WORKDIR /app
 ---> Running in xyz789uvw012
...
Successfully tagged myapp:1.0.0
```

**Example:**
```bash
curl -X POST http://localhost:8089/images/build \
  -H "Content-Type: application/json" \
  -d '{"context_path":"D:\\Code\\myapp","dockerfile":"Dockerfile.dev","tag":"myapp:dev"}'
```

### Volumes
```bash
curl http://localhost:8089/volumes
```

---

### Inspect Volume

**GET** `/volumes/inspect?name=<volume_name>`

Gets detailed information about a volume.

**Parameters:**
- `name` (required): Volume name

**Response:**
```json
{
  "name": "database_data",
  "driver": "local",
  "mountpoint": "/var/lib/docker/volumes/database_data/_data",
  "scope": "local",
  "labels": {"backup": "daily"},
  "created_at": "2024-01-10T10:30:00Z",
  "in_use": true
}
```

**Example:**
```bash
curl "http://localhost:8089/volumes/inspect?name=database_data"
```

---

### Create Volume

**POST** `/volumes/create`

Creates a new volume.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "new_volume",
  "driver": "local",
  "labels": {
    "backup": "daily",
    "environment": "production"
  },
  "options": {
    "type": "tmpfs",
    "size": "100m"
  }
}
```

**Parameters:**
- `name` (required): Unique volume name
- `driver` (optional): Volume driver (default: "local")
- `labels` (optional): Metadata labels
- `options` (optional): Driver-specific options

**Response:**
```json
{
  "status": "created",
  "name": "new_volume"
}
```

**Example:**
```bash
curl -X POST http://localhost:8089/volumes/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "app_data",
    "driver": "local"
  }'
```

---

### Remove Volume

**DELETE** `/volumes/remove?name=<volume_name>`

Removes a volume.

**Parameters:**
- `name` (required): Volume name

**Response:**
```json
{
  "status": "removed",
  "name": "database_data"
}
```

**Example:**
```bash
curl -X DELETE "http://localhost:8089/volumes/remove?name=database_data"
```

---

### Prune Unused Volumes

**POST** `/volumes/prune`

Removes all unused volumes.

**Response:**
```json
{
  "volumes_deleted": ["unused_vol1", "unused_vol2"],
  "space_reclaimed": 2147483648
}
```

**Example:**
```bash
curl -X POST http://localhost:8089/volumes/prune
```

### Networks
```bash
curl http://localhost:8089/networks
```

---

### Create Network

**POST** `/networks/create`

Creates a new network.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "app_network",
  "driver": "bridge",
  "options": {
    "com.docker.network.bridge.name": "br_app"
  }
}
```

**Parameters:**
- `name` (required): Network name
- `driver` (optional): Driver type (bridge, overlay, macvlan, ipvlan, host)
- `options` (optional): Driver-specific options

**Response:**
```json
{
  "id": "abc123def456",
  "name": "app_network",
  "driver": "bridge"
}
```

**Example:**
```bash
curl -X POST http://localhost:8089/networks/create \
  -H "Content-Type: application/json" \
  -d '{"name":"scout-net","driver":"bridge","options":{}}'
curl -X DELETE "http://localhost:8089/networks/remove?id=scout-net"
```

### Statistics & Monitoring
```bash
curl -X POST http://localhost:8089/events/publish \
  -H "Content-Type: application/json" \
  -d '{
    "type": "custom.alert",
    "message": "Test alert"
  }'
```

---

## Statistics

### Get System Statistics

**GET** `/stats`

Returns current system statistics snapshot.

**Response:**
```json
{
  "summary": {
    "active_containers": 5,
    "avg_cpu": "25.4%",
    "avg_memory": "2.1GB",
    "avg_net_io": "15.2 Mbps",
    "avg_disk_io": "8.3 MB/s"
  },
  "containers": [
    {
      "id": "abc123def456",
      "name": "web-server",
      "state": "running",
      "image": "nginx:latest",
      "cpu": "15.2%",
      "memory": "512MB",
      "net_io": "5.2 Mbps",
      "disk_io": "2.1 MB/s",
      "last_used": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Update Frequency:** Real-time via WebSocket (every 1 second)

---

## WebSocket Endpoints

### Stats WebSocket

**WS** `/stats`

Real-time Docker container statistics streamed every second.

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8089/stats');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Stats:', data);
};
```

**Message Format:**
```json
{
  "type": "stats",
  "timestamp": 1234567890,
  "stats": {
    "container_id_1": {
      "id": "abc123",
      "name": "/web-server",
      "cpu_percent": 12.5,
      "memory_bytes": 104857600,
      "memory_limit": 2147483648,
      "net_rx_bps": 5000,
      "net_tx_bps": 3000,
      "disk_read_bps": 1024,
      "disk_write_bps": 2048
    }
  }
}
```

**Features:**
- Auto-updates every 1 second
- All running containers included
- Real-time CPU, memory, network, and disk stats

---

### Terminal WebSocket

**WS** `/terminal`

Execute shell commands and stream output in real-time.

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8089/terminal');

// Send command to execute
ws.send(JSON.stringify({
  type: 'start',
  command: 'docker ps -a'
}));

// Receive output
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  switch(data.type) {
    case 'started':
      console.log('Command started');
      break;
    case 'stdout':
      console.log('Output:', data.data);
      break;
    case 'stderr':
      console.error('Error:', data.data);
      break;
    case 'exit':
      console.log('Exited with code:', data.code);
      break;
  }
};

// Terminate running command
ws.send(JSON.stringify({ type: 'exitTerm' }));
```

**Request Messages:**

1. **Start Command:**
```json
{
  "type": "start",
  "command": "your shell command here"
}
```

2. **Terminate Command:**
```json
{
  "type": "exitTerm"
}
```

**Response Messages:**

1. **Command Started:**
```json
{
  "type": "started",
  "message": "command started"
}
```

2. **Standard Output (per line):**
```json
{
  "type": "stdout",
  "data": "output line"
}
```

3. **Standard Error (per line):**
```json
{
  "type": "stderr",
  "data": "error line"
}
```

4. **Command Exit:**
```json
{
  "type": "exit",
  "code": 0,
  "message": "command completed successfully"
}
```

5. **Command Terminated:**
```json
{
  "type": "terminated",
  "message": "command terminated"
}
```

6. **Error:**
```json
{
  "type": "error",
  "message": "error description"
}
```

**Features:**
- Real-time output streaming (line by line)
- Separate stdout and stderr streams
- Graceful termination support
- Exit code reporting

---

## Complete Route Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/system/health` | Health check |
| POST | `/system/nuke` | System cleanup |
| GET | `/containers` | List containers |
| POST | `/containers/start` | Start container |
| POST | `/containers/stop` | Stop container |
| POST | `/containers/pause` | Pause container |
| POST | `/containers/unpause` | Unpause container |
| DELETE | `/containers/remove` | Remove container |
| POST | `/containers/create` | Create container |
| GET | `/images` | List images |
| DELETE | `/images/remove` | Remove image |
| POST | `/images/build` | Build image |
| GET | `/volumes` | List volumes |
| GET | `/volumes/inspect` | Inspect volume |
| POST | `/volumes/create` | Create volume |
| DELETE | `/volumes/remove` | Remove volume |
| POST | `/volumes/prune` | Prune unused volumes |
| GET | `/networks` | List networks |
| POST | `/networks/create` | Create network |
| DELETE | `/networks/remove` | Remove network |
<<<<<<< HEAD
| GET | `/events` | WebSocket stream |
| POST | `/events/publish` | Publish event |
| GET | `/terminal` | WebSocket terminal |
| GET | `/stats` | System statistics |
=======
| WS | `/stats` | Real-time stats stream |
| WS | `/terminal` | Command execution stream |
>>>>>>> 22f7920 (Updated WebSocket)
| GET | `/` | Dashboard UI |
| GET | `/echo` | Terminal UI |
| GET | `/index` | Terminal UI (alias) |

---

## Data Models

### Container

```json
{
  "id": "string",
  "name": "string",
  "state": "string",
  "image": "string",
  "ports": "string",
  "labels": {"key": "value"},
  "cpu": "string",
  "memory": "string",
  "net_io": "string",
  "disk_io": "string",
  "last_used": "string"
}
```

### Image

```json
{
  "id": "string",
  "repo_tags": ["string"],
  "repo_digests": ["string"],
  "size": "number",
  "created": "string"
}
```

### Volume

```json
{
  "name": "string",
  "driver": "string",
  "mountpoint": "string",
  "scope": "string",
  "labels": {"key": "value"},
  "created_at": "string",
  "in_use": "boolean"
}
```

### Network

```json
{
  "id": "string",
  "name": "string",
  "driver": "string",
  "scope": "string",
  "containers": "number",
  "created": "string"
}
```

---

## Error Handling

### Standard Error Response

```json
{
  "error": "descriptive error message"
}
```

### HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | OK - Successful request |
| 201 | Created - Resource created |
| 202 | Accepted - Event published |
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found - Resource not found |
| 405 | Method Not Allowed - Wrong HTTP method |
| 409 | Conflict - Resource in use |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Server busy |

### Common Errors

| Error | Solution |
|-------|----------|
| `container id required` | Provide container ID in query parameter |
| `volume is attached to a container` | Disconnect container before removing |
| `network is connected to running containers` | Stop/disconnect containers first |
| `image is in use by containers` | Remove containers using image first |
| `streaming unsupported` | Use HTTP/1.1 or HTTP/2 |

---

## Examples

### Complete Workflow

```bash
# 1. Check health
curl http://localhost:8089/system/health

# 2. List containers
curl http://localhost:8089/containers

# 3. List images
curl http://localhost:8089/images

# 4. Create a container
curl -X POST http://localhost:8089/containers/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-nginx",
    "image": "nginx:latest",
    "port_bindings": {"80/tcp": "8080"},
    "restart_policy": "unless-stopped"
  }'

# 5. Start the container
curl -X POST "http://localhost:8089/containers/start?id=my-nginx"

# 6. Build new image
curl -X POST http://localhost:8089/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "context_path": "/path/to/app",
    "tag": "myapp:1.0.0"
  }'

# 5. Create volume
curl -X POST http://localhost:8089/volumes/create \
  -H "Content-Type: application/json" \
  -d '{"name": "app_data", "driver": "local"}'

# 6. Create network
curl -X POST http://localhost:8089/networks/create \
  -H "Content-Type: application/json" \
  -d '{"name": "app_network", "driver": "bridge"}'

# 7. Get statistics
curl http://localhost:8089/stats

# 8. Stop container
curl -X POST "http://localhost:8089/containers/stop?id=mycontainer"

# 9. Remove container
curl -X DELETE "http://localhost:8089/containers/remove?id=mycontainer"
```

### Using with JavaScript/Node.js

```javascript
// Simple async function to fetch containers
async function listContainers() {
  const response = await fetch('http://localhost:8089/containers');
  const data = await response.json();
  console.log(data);
}

// Start a container
async function startContainer(containerId) {
  const response = await fetch(`http://localhost:8089/containers/start?id=${containerId}`, {
    method: 'POST'
  });
  const data = await response.json();
  console.log(data);
}

// Listen to real-time stats
function listenToStats() {
  const ws = new WebSocket('ws://localhost:8089/stats');
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Stats update:', data.stats);
  };
}

// Execute command and stream output
function executeCommand(command) {
  const ws = new WebSocket('ws://localhost:8089/terminal');
  
  ws.onopen = () => {
    ws.send(JSON.stringify({ type: 'start', command }));
  };
  
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    
    if (msg.type === 'stdout') {
      console.log(msg.data);
    } else if (msg.type === 'stderr') {
      console.error(msg.data);
    } else if (msg.type === 'exit') {
      console.log('Exit code:', msg.code);
      ws.close();
    }
  };
  
  // To terminate: ws.send(JSON.stringify({ type: 'exitTerm' }));
}
```

---

## Troubleshooting

### Connection Refused

```bash
# Verify server is running
docker compose ps

# Check logs
docker compose logs -f docker_api_server

# Restart if needed
docker compose restart docker_api_server
```

### Permission Denied

```bash
# Ensure Docker socket has correct permissions
sudo chmod 666 /var/run/docker.sock

# Or add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Port Already in Use

```bash
# Use different port
export DOCKER_SCOUT_ADDR=:9000
docker compose up -d
```

### Build Timeout

```bash
# Check Docker daemon logs
docker logs docker_api_server

# Increase timeout or check network connectivity
```

---

## Security Notes

⚠️ **Current Implementation:**
- No authentication mechanism
- CORS enabled for all origins
- Accessible to anyone on the network

✅ **Production Recommendations:**
1. Deploy behind authentication gateway (Kong, Nginx + OAuth2)
2. Use TLS/HTTPS
3. Implement network policies
4. Restrict CORS to specific origins
5. Add request logging and monitoring
6. Run with least privilege user
7. Use Docker socket permissions for access control

---

## Performance

### Timeouts
- Container operations: 10 seconds
- Network operations: 10-15 seconds
- Volume operations: 10 seconds
- Image builds: Streaming (no timeout)

### Resource Usage
- Memory: ~50MB at startup
- CPU: <5% idle
- Scales with number of active containers/images

---

## Support

For detailed API documentation, see: [COMPLETE_API_DOCS.md](COMPLETE_API_DOCS.md)

---

## License

See LICENSE file in repository.

---

**Happy container managing! 🐳**
