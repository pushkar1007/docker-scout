Docker Scout Server

This server exposes a REST API + SSE stream for Docker resources.

**Default server URL:** `http://localhost:8089`

## Run (Linux/macOS)
From `server/`:

```bash
go run ./cmd/scout
```

CLI mode:

```bash
go run ./cmd/scout --cli
```

## Run (Windows)
From `server\` (Git Bash):

```bash
go run ./cmd/scout
```

CLI mode:

```bash
go run ./cmd/scout --cli
```

If Docker Desktop is running, the client will use the named pipe by default.

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
curl http://localhost:8089/system/health
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
curl -X DELETE "http://localhost:8089/images/remove?name=<image_name>&force=true"
curl -X POST "http://localhost:8089/images/build" \
  -H "Content-Type: application/json" \
  -d '{"context_path":"D:\\Code\\myapp","dockerfile":"Dockerfile.dev","tag":"myapp:dev"}'
```

### Volumes
```bash
curl http://localhost:8089/volumes
curl "http://localhost:8089/volumes/inspect?name=<volume_name>"
curl -X POST http://localhost:8089/volumes/create \
  -H "Content-Type: application/json" \
  -d '{"name":"volume_name","driver":"local"}'
curl -X DELETE "http://localhost:8089/volumes/remove?name=<volume_name>"
curl -X POST http://localhost:8089/volumes/prune
```

### Networks
```bash
curl http://localhost:8089/networks
curl -X POST http://localhost:8089/networks/create \
  -H "Content-Type: application/json" \
  -d '{"name":"scout-net","driver":"bridge","options":{}}'
curl -X DELETE "http://localhost:8089/networks/remove?id=scout-net"
```

### Statistics & Monitoring
```bash
curl http://localhost:8089/stats
curl -N http://localhost:8089/events
curl -X POST http://localhost:8089/events/publish \
  -H "Content-Type: application/json" \
  -d '{"event":"custom_event","data":"hello"}'
```
