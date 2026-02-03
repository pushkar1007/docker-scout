Docker Scout Server

This server exposes a REST API + SSE stream for Docker resources.

**Default server URL:** `http://localhost:8089`

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

**Update Frequency:** Every 5 seconds

**Example:**
```bash
curl http://localhost:8089/stats
```

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
| GET | `/events` | SSE stream |
| POST | `/events/publish` | Publish event |
| GET | `/stats` | System statistics |
| GET | `/` | Dashboard UI |

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

// Listen to events
function listenToEvents() {
  const eventSource = new EventSource('http://localhost:8089/events');
  eventSource.onmessage = (event) => {
    console.log('Event:', JSON.parse(event.data));
  };
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