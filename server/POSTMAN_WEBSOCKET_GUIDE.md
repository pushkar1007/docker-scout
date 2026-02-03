# WebSocket Testing Guide for Postman

## 1. Stats WebSocket - Real-time Docker Statistics

**WebSocket URL:**
```
ws://localhost:8089/stats
```

### How to Test in Postman:

1. Create a new **WebSocket Request**
2. Enter URL: `ws://localhost:8089/stats`
3. Click **Connect**
4. You will automatically receive stats updates every 1 second

### Expected Response (every second):

```json
{
  "type": "stats",
  "timestamp": 1738627200,
  "stats": {
    "abc123def456": {
      "id": "abc123def456",
      "name": "/my-container",
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

### Notes:
- No messages need to be sent
- Stats auto-update every second
- Includes all running containers
- Connection stays open until you disconnect

---

## 2. Terminal WebSocket - Execute Commands

**WebSocket URL:**
```
ws://localhost:8089/terminal
```

### How to Test in Postman:

1. Create a new **WebSocket Request**
2. Enter URL: `ws://localhost:8089/terminal`
3. Click **Connect**
4. Send messages to execute commands

---

### Message 1: Start a Command

**Send this message:**
```json
{
  "type": "start",
  "command": "docker ps -a"
}
```

**You will receive:**

1. Started confirmation:
```json
{
  "type": "started",
  "message": "command started"
}
```

2. Output lines (stdout):
```json
{
  "type": "stdout",
  "data": "CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS"
}
```

3. Error lines (stderr) if any:
```json
{
  "type": "stderr",
  "data": "error message here"
}
```

4. Exit status when complete:
```json
{
  "type": "exit",
  "code": 0,
  "message": "command completed successfully"
}
```

---

### Message 2: Terminate Running Command

**Send this message:**
```json
{
  "type": "exitTerm"
}
```

**You will receive:**
```json
{
  "type": "terminated",
  "message": "command terminated"
}
```

---

## Example Commands to Test

### List Docker Containers
```json
{
  "type": "start",
  "command": "docker ps -a"
}
```

### Check System Info
```json
{
  "type": "start",
  "command": "docker info"
}
```

### Long-running Command (watch Docker stats)
```json
{
  "type": "start",
  "command": "docker stats --no-stream"
}
```

To stop it:
```json
{
  "type": "exitTerm"
}
```

### Ping Test (continuous output)
```json
{
  "type": "start",
  "command": "ping -c 5 google.com"
}
```

---

## Error Handling

If you send an invalid message, you'll receive:
```json
{
  "type": "error",
  "message": "error description here"
}
```

Common errors:
- `"invalid request format"` - JSON parsing failed
- `"command is required"` - Missing command in start request
- `"unknown request type: xyz"` - Invalid type field

---

## Quick Reference

| Endpoint | Purpose | Auto-sends? | Need to send? |
|----------|---------|-------------|---------------|
| `/stats` | Real-time Docker stats | ✅ Yes (every 1s) | ❌ No |
| `/terminal` | Execute shell commands | ❌ No | ✅ Yes |

---

## Postman Tips

1. **Multiple tabs**: Open both WebSockets in separate tabs
2. **Save messages**: Save common commands as examples
3. **Message history**: Postman shows all sent/received messages
4. **Connection status**: Watch for "Connected" indicator
5. **Disconnect**: Click "Disconnect" when done testing

---

## Testing Workflow

### Stats WebSocket:
```
1. Connect → 2. Watch auto-updates → 3. Disconnect when done
```

### Terminal WebSocket:
```
1. Connect → 2. Send start command → 3. Watch output → 4. (Optional) Send exitTerm → 5. Disconnect
```

---

## ⚠️ Troubleshooting

### Error: connect ECONNREFUSED 127.0.0.1:8089

**This means the server is not running!**

#### Start the server:

**Option 1: Docker Compose (Recommended)**
```bash
cd /home/atithisingh/Desktop/desktop-tech/subdomain-proxy/server
docker compose up -d
```

**Option 2: Run directly with Go**
```bash
cd /home/atithisingh/Desktop/desktop-tech/subdomain-proxy/server
go run cmd/scout/main.go
```

#### Verify server is running:
```bash
curl http://localhost:8089/system/health
```

Expected response:
```json
{"status":"ok"}
```

#### Check server logs:
```bash
# Docker logs
docker compose logs -f

# Or if running with Go
# Check terminal output
```

#### Common Issues:

1. **Port 8089 already in use**
   ```bash
   # Find what's using port 8089
   lsof -i :8089
   
   # Or use different port
   export DOCKER_SCOUT_ADDR=:9000
   go run cmd/scout/main.go
   ```

2. **Docker not running**
   ```bash
   sudo systemctl start docker
   # Or start Docker Desktop
   ```

3. **Permission denied on Docker socket**
   ```bash
   sudo chmod 666 /var/run/docker.sock
   ```

---

**Server must be running at:** `http://localhost:8089`

**Before testing WebSockets, verify with:**
```bash
curl http://localhost:8089/system/health
```

If you get `{"status":"ok"}`, you're ready to test WebSockets in Postman! ✅
