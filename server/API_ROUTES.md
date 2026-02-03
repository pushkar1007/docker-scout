# Docker Scout API Routes

## Overview
A complete REST API server for managing Docker containers, images, volumes, and system resources using the Docker API.

## Running the Server
```bash
docker compose up -d
# Server runs on http://localhost:8089
```

## API Endpoints

### System Management

#### Health Check
```bash
GET /system/health
```
Response: `{"status":"healthy","version":"1.0.0"}`

#### System Cleanup (Nuclear Option)
```bash
POST /system/nuke?confirm=true
```
Warning: Removes all containers, images, and unused volumes.
Response: `{"status":"system cleaned"}`

---

### Containers

#### List All Containers
```bash
GET /containers
```
Returns list of all running containers with details.

#### Start Container
```bash
POST /containers/start?id=<container_id>
```
Response: `{"status":"started","id":"<container_id>"}`

#### Stop Container
```bash
POST /containers/stop?id=<container_id>
```
Response: `{"status":"stopped","id":"<container_id>"}`

#### Pause Container
```bash
POST /containers/pause?id=<container_id>
```
Response: `{"status":"paused","id":"<container_id>"}`

#### Unpause Container
```bash
POST /containers/unpause?id=<container_id>
```
Response: `{"status":"unpaused","id":"<container_id>"}`

#### Remove Container
```bash
DELETE /containers/remove?id=<container_id>
```
Response: `{"status":"removed","id":"<container_id>"}`

---

### Images

#### List All Images
```bash
GET /images
```
Returns list of all images with metadata.

#### Remove Image
```bash
DELETE /images/remove?name=<image_name>&force=true
```
Response: `{"status":"removed","image":"<image_name>"}`

#### Build Image
```bash
POST /images/build
Content-Type: application/json

{
  "context_path": "/path/to/dockerfile",
  "tag": "image:tag"
}
```
Response: Build output stream

---

### Volumes

#### List All Volumes
```bash
GET /volumes
```
Returns list of all volumes.

#### Inspect Volume
```bash
GET /volumes/inspect?name=<volume_name>
```
Returns detailed volume information.

#### Create Volume
```bash
POST /volumes/create
Content-Type: application/json

{
  "name": "volume_name",
  "driver": "local"
}
```
Response: `{"status":"created","name":"volume_name"}`

#### Remove Volume
```bash
DELETE /volumes/remove?name=<volume_name>
```
Response: `{"status":"removed","name":"volume_name"}`

#### Prune Unused Volumes
```bash
POST /volumes/prune
```
Response: Report of pruned volumes.

---

### Statistics & Monitoring

#### Get Container Statistics
```bash
GET /stats
```
Returns:
```json
{
  "summary": {
    "active_containers": 2,
    "avg_cpu": "0.04%",
    "avg_memory": "12.10MiB",
    "avg_net_io": "0B/s",
    "avg_disk_io": "0B/s"
  },
  "containers": [...]
}
```

#### Server-Sent Events Stream
```bash
GET /events
```
Real-time event stream of Docker events.

#### Publish Custom Event
```bash
POST /events/publish
Content-Type: application/json

{"event": "custom_event", "data": "event_data"}
```

---

## Response Format

All successful responses use standard HTTP status codes:
- **200 OK**: Successful GET/POST
- **201 Created**: Resource created
- **204 No Content**: Successful action with no response body
- **400 Bad Request**: Invalid parameters
- **409 Conflict**: Resource conflict (e.g., volume in use)
- **500 Internal Server Error**: Server error
- **503 Service Unavailable**: Server too busy

All error responses include JSON:
```json
{"error": "error message"}
```

## Architecture

- **Language**: Go
- **Server**: net/http (standard library)
- **Docker Client**: moby/moby/client
- **Base Image**: distroless/base-debian12 (minimal)
- **Port**: 8089
- **Docker Socket**: Mounted from host via `/var/run/docker.sock`

## File Structure

```
cmd/scout/main.go              - Entry point
internal/api/
  - routes.go                 - Route registration
  - containers.go             - Container endpoints
  - images.go                 - Image endpoints
  - volumes.go                - Volume endpoints
  - stats.go                  - Statistics endpoint
  - events.go                 - Event streaming
  - system.go                 - System management
internal/docker/
  - client.go                 - Docker client initialization
  - containers.go             - Container operations
  - images.go                 - Image operations
  - volumes.go                - Volume operations
  - stats.go                  - Statistics collection
internal/app/
  - server.go                 - Server startup
  - cli.go                    - CLI interface
  - lifecycle.go              - Stats updater goroutine
internal/state/
  - cache.go                  - In-memory cache
  - broadcaster.go            - Event broadcasting
  - updater.go                - Cache updater
internal/model/
  - container.go              - Container models
  - responses.go              - API response models
  - stats.go                  - Statistics models
```

## Features

✅ Container lifecycle management (start, stop, pause, unpause)
✅ Image listing and removal
✅ Volume management (create, remove, inspect)
✅ Real-time statistics and monitoring
✅ Server-sent events (SSE) for live updates
✅ System cleanup utilities
✅ Proper JSON responses with CORS headers
✅ Health check endpoint
✅ Docker socket integration
✅ Asynchronous stats updates
