# Docker Scout - Docker Management API

A comprehensive REST API server for managing Docker containers, images, volumes, and networks with real-time monitoring and streaming capabilities.

## Features

✨ **Complete Docker Management**
- Container lifecycle management (start, stop, pause, unpause, remove)
- Image management (list, build, remove)
- Volume management (create, list, inspect, remove, prune)
- Network management (create, list, remove)

📊 **Real-Time Monitoring**
- Live Server-Sent Events (SSE) streaming
- System statistics and resource monitoring
- Container performance metrics (CPU, memory, network I/O, disk I/O)
- Dashboard UI for visualization

🚀 **Production Ready**
- Cross-Origin Resource Sharing (CORS) enabled
- Comprehensive error handling
- Request timeouts and context management
- Streaming support for image builds and events

---

## Quick Start

### Prerequisites
- Docker installed and running
- Docker Compose installed
- Linux/macOS or WSL2 on Windows

### Installation & Running

```bash
# Clone the repository
cd server

# Start the server
docker compose up -d

# Verify it's running
curl http://localhost:3000/system/health
```

**Server is now running at:** `http://localhost:3000`

### Stop the Server

```bash
docker compose down
```

### View Logs

```bash
docker compose logs -f docker_api_server
```

---

## Configuration

### Environment Variables

```bash
# Set custom port (default: 3000)
export DOCKER_SCOUT_ADDR=:9000
docker compose up -d
```

**Environment File (.env):**
```env
DOCKER_SCOUT_ADDR=:3000
```

Then run:
```bash
docker compose up -d
```

---

## API Documentation

### Base URL

```
http://localhost:3000
```

### Response Format

All responses are in JSON format with CORS headers enabled.

---

## System Endpoints

### Health Check

**GET** `/system/health`

Health status and version check.

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

**Example:**
```bash
curl http://localhost:3000/system/health
```

---

### System Cleanup

**POST** `/system/nuke?confirm=true`

⚠️ **DESTRUCTIVE:** Removes all containers, images, and unused volumes.

**Parameters:**
- `confirm` (required): Must be `"true"`

**Response:**
```json
{
  "status": "system cleaned"
}
```

**Example:**
```bash
curl -X POST "http://localhost:3000/system/nuke?confirm=true"
```

---

## Container Endpoints

### List All Containers

**GET** `/containers`

Lists all containers (running and stopped).

**Response:**
```json
{
  "items": [
    {
      "id": "abc123def456",
      "name": "web-server",
      "state": "running",
      "image": "nginx:latest",
      "ports": "80:8080",
      "labels": {"app": "web"},
      "cpu": "15.2%",
      "memory": "512MB",
      "net_io": "5.2 Mbps",
      "disk_io": "2.1 MB/s",
      "last_used": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:3000/containers
```

---

### Start Container

**POST** `/containers/start?id=<container_id>`

Starts a stopped container.

**Parameters:**
- `id` (required): Container ID or name

**Response:**
```json
{
  "status": "started",
  "id": "web-server"
}
```

**Example:**
```bash
curl -X POST "http://localhost:3000/containers/start?id=web-server"
```

---

### Stop Container

**POST** `/containers/stop?id=<container_id>`

Gracefully stops a running container.

**Parameters:**
- `id` (required): Container ID or name

**Response:**
```json
{
  "status": "stopped",
  "id": "web-server"
}
```

**Example:**
```bash
curl -X POST "http://localhost:3000/containers/stop?id=web-server"
```

---

### Pause Container

**POST** `/containers/pause?id=<container_id>`

Pauses a running container (freezes processes).

**Parameters:**
- `id` (required): Container ID or name

**Response:**
```json
{
  "status": "paused",
  "id": "web-server"
}
```

**Example:**
```bash
curl -X POST "http://localhost:3000/containers/pause?id=web-server"
```

---

### Unpause Container

**POST** `/containers/unpause?id=<container_id>`

Resumes a paused container.

**Parameters:**
- `id` (required): Container ID or name

**Response:**
```json
{
  "status": "unpaused",
  "id": "web-server"
}
```

**Example:**
```bash
curl -X POST "http://localhost:3000/containers/unpause?id=web-server"
```

---

### Remove Container

**DELETE** `/containers/remove?id=<container_id>`

Removes a container permanently.

**Parameters:**
- `id` (required): Container ID or name

**Response:**
```json
{
  "status": "removed",
  "id": "web-server"
}
```

**Example:**
```bash
curl -X DELETE "http://localhost:3000/containers/remove?id=web-server"
```

---

### Create Container

**POST** `/containers/create`

Creates a new container from an image.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "my-app",
  "image": "nginx:latest",
  "cmd": ["/bin/sh", "-c", "nginx -g 'daemon off;'"],
  "env": ["NODE_ENV=production", "PORT=3000"],
  "labels": {
    "app": "web",
    "environment": "production"
  },
  "exposed_ports": {
    "80/tcp": {},
    "443/tcp": {}
  },
  "port_bindings": {
    "80/tcp": "8080",
    "443/tcp": "8443"
  },
  "volumes": [
    "/host/path:/container/path",
    "volume-name:/data"
  ],
  "network_mode": "bridge",
  "restart_policy": "unless-stopped"
}
```

**Parameters:**
- `name` (optional): Container name (auto-generated if not provided)
- `image` (required): Image name and tag
- `cmd` (optional): Command to run in container
- `env` (optional): Environment variables (array of "KEY=VALUE")
- `labels` (optional): Container labels
- `exposed_ports` (optional): Ports to expose (format: "port/protocol")
- `port_bindings` (optional): Host port mappings (format: "container_port/protocol": "host_port")
- `volumes` (optional): Volume mounts (format: "host:container" or "volume:container")
- `network_mode` (optional): Network mode (bridge, host, none, container:<name|id>)
- `restart_policy` (optional): Restart policy (no, always, unless-stopped, on-failure)

**Response:**
```json
{
  "status": "created",
  "id": "abc123def456",
  "name": "my-app"
}
```

**Example:**
```bash
# Simple container
curl -X POST http://localhost:3000/containers/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web-server",
    "image": "nginx:latest"
  }'

# With port mappings and environment
curl -X POST http://localhost:3000/containers/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-app",
    "image": "node:18-alpine",
    "env": ["NODE_ENV=production", "PORT=3000"],
    "port_bindings": {"3000/tcp": "8080"},
    "restart_policy": "unless-stopped"
  }'

# With volumes and labels
curl -X POST http://localhost:3000/containers/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "database",
    "image": "postgres:14",
    "env": ["POSTGRES_PASSWORD=secret"],
    "volumes": ["db-data:/var/lib/postgresql/data"],
    "labels": {"tier": "database"},
    "restart_policy": "always"
  }'
```

**Note:** Add `?start=true` to automatically start the container after creation. All fields (env, labels, volumes, cmd, etc.) are applied when creating the container.

**Example Usage (Create & Start in One Request):**
```bash
# Simple container - create and start
curl -X POST "http://localhost:3000/containers/create?start=true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "web-server",
    "image": "nginx:latest"
  }'


# With all options - PostgreSQL with env, volumes, labels, restart policy
curl -X POST "http://localhost:3000/containers/create?start=true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "database",
    "image": "postgres:14",
    "env": ["POSTGRES_PASSWORD=secret"],
    "volumes": ["db-data:/var/lib/postgresql/data"],
    "labels": {"tier": "database", "env": "production"},
    "restart_policy": "unless-stopped"
  }'

# Verify it's running
curl http://localhost:3000/containers | jq '.items[] | {name, state, image, labels}'
```

---

## Image Endpoints

### List All Images

**GET** `/images`

Lists all available images.

**Response:**
```json
{
  "items": [
    {
      "id": "sha256:abc123def456",
      "repo_tags": ["nginx:latest", "nginx:1.25"],
      "repo_digests": ["nginx@sha256:..."],
      "size": 187395149,
      "created": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:3000/images
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
curl -X DELETE "http://localhost:3000/images/remove?name=nginx:latest"
curl -X DELETE "http://localhost:3000/images/remove?name=nginx:latest&force=true"
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
curl -X POST http://localhost:3000/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "context_path": "/path/to/app",
    "tag": "myapp:1.0.0"
  }'
```

---

## Volume Endpoints

### List All Volumes

**GET** `/volumes`

Lists all volumes.

**Response:**
```json
{
  "volumes": [
    {
      "name": "database_data",
      "driver": "local",
      "mountpoint": "/var/lib/docker/volumes/database_data/_data",
      "scope": "local",
      "labels": {"backup": "daily"},
      "created_at": "2024-01-10T10:30:00Z",
      "in_use": true
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:3000/volumes
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
curl "http://localhost:3000/volumes/inspect?name=database_data"
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
curl -X POST http://localhost:3000/volumes/create \
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
curl -X DELETE "http://localhost:3000/volumes/remove?name=database_data"
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
curl -X POST http://localhost:3000/volumes/prune
```

---

## Network Endpoints

### List All Networks

**GET** `/networks`

Lists all Docker networks.

**Response:**
```json
[
  {
    "id": "abc123def456",
    "name": "bridge",
    "driver": "bridge",
    "scope": "local",
    "containers": 2,
    "created": "2024-01-01T00:00:00Z"
  }
]
```

**Example:**
```bash
curl http://localhost:3000/networks
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
curl -X POST http://localhost:3000/networks/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "app_network",
    "driver": "bridge"
  }'
```

---

### Remove Network

**DELETE** `/networks/remove?id=<network_id>` or `?name=<network_name>`

Removes a network.

**Parameters:**
- `id` (optional): Network ID
- `name` (optional): Network name

**Response:**
```json
{
  "status": "removed",
  "id": "abc123def456"
}
```

**Example:**
```bash
curl -X DELETE "http://localhost:3000/networks/remove?name=app_network"
curl -X DELETE "http://localhost:3000/networks/remove?id=abc123def456"
```

---

## Events & Streaming

### Server-Sent Events Stream

**GET** `/events`

Real-time event stream (Server-Sent Events).

**Response Headers:**
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

**Response (Streaming):**
```
data: {"type":"container.started","container_id":"abc123","name":"web-server","timestamp":"2024-01-15T10:30:00Z"}

data: {"type":"container.stopped","container_id":"xyz789","name":"database","timestamp":"2024-01-15T10:31:00Z"}

data: {"type":"stats.updated","containers":5,"cpu":"25.4%","memory":"2.1GB"}
```

**Example:**
```bash
# Watch live events
curl http://localhost:3000/events

# Or with jq for formatting
curl http://localhost:3000/events | jq .
```

**JavaScript Client:**
```javascript
const eventSource = new EventSource('http://localhost:3000/events');

eventSource.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
});

eventSource.onerror = (error) => {
  console.error('Connection error:', error);
  eventSource.close();
};
```

---

### Publish Event

**POST** `/events/publish`

Publish a custom event to all connected clients.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "type": "custom.alert",
  "severity": "warning",
  "message": "High memory usage detected",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Response:**
```json
{
  "status": "published"
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/events/publish \
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
curl http://localhost:3000/stats
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
curl http://localhost:3000/system/health

# 2. List containers
curl http://localhost:3000/containers

# 3. List images
curl http://localhost:3000/images

# 4. Create a container
curl -X POST http://localhost:3000/containers/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-nginx",
    "image": "nginx:latest",
    "port_bindings": {"80/tcp": "8080"},
    "restart_policy": "unless-stopped"
  }'

# 5. Start the container
curl -X POST "http://localhost:3000/containers/start?id=my-nginx"

# 6. Build new image
curl -X POST http://localhost:3000/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "context_path": "/path/to/app",
    "tag": "myapp:1.0.0"
  }'

# 5. Create volume
curl -X POST http://localhost:3000/volumes/create \
  -H "Content-Type: application/json" \
  -d '{"name": "app_data", "driver": "local"}'

# 6. Create network
curl -X POST http://localhost:3000/networks/create \
  -H "Content-Type: application/json" \
  -d '{"name": "app_network", "driver": "bridge"}'

# 7. Get statistics
curl http://localhost:3000/stats

# 8. Stop container
curl -X POST "http://localhost:3000/containers/stop?id=mycontainer"

# 9. Remove container
curl -X DELETE "http://localhost:3000/containers/remove?id=mycontainer"
```

### Using with JavaScript/Node.js

```javascript
// Simple async function to fetch containers
async function listContainers() {
  const response = await fetch('http://localhost:3000/containers');
  const data = await response.json();
  console.log(data);
}

// Start a container
async function startContainer(containerId) {
  const response = await fetch(`http://localhost:3000/containers/start?id=${containerId}`, {
    method: 'POST'
  });
  const data = await response.json();
  console.log(data);
}

// Listen to events
function listenToEvents() {
  const eventSource = new EventSource('http://localhost:3000/events');
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