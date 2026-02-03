package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type terminalStartRequest struct {
	Type        string   `json:"type"`
	ContainerID string   `json:"container_id"`
	Cmd         []string `json:"cmd"`
	TTY         *bool    `json:"tty"`
	Env         []string `json:"env"`
	Workdir     string   `json:"workdir"`
	User        string   `json:"user"`
	Cols        uint     `json:"cols"`
	Rows        uint     `json:"rows"`
}

type terminalResizeRequest struct {
	Type string `json:"type"`
	Cols uint   `json:"cols"`
	Rows uint   `json:"rows"`
}

var terminalUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func registerTerminal(mux *http.ServeMux, _ Deps) {
	mux.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conn, err := terminalUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		writeMu := sync.Mutex{}
		writeText := func(payload []byte) {
			writeMu.Lock()
			_ = conn.WriteMessage(websocket.TextMessage, payload)
			writeMu.Unlock()
		}
		writeError := func(msg string) {
			writeText([]byte("error: " + msg))
		}

		_, first, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var start terminalStartRequest
		if err := json.Unmarshal(first, &start); err != nil || start.Type != "start" {
			writeError("invalid start payload")
			return
		}

		inputBuf := make([]byte, 0, 4096)
		const maxCommandBuffer = 64 * 1024

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if mt == websocket.TextMessage {
				trimmed := strings.TrimSpace(string(msg))
				if strings.HasPrefix(trimmed, "{") {
					var resize terminalResizeRequest
					if json.Unmarshal(msg, &resize) == nil && resize.Type == "resize" {
						if resize.Cols > 0 && resize.Rows > 0 {
							start.Cols = resize.Cols
							start.Rows = resize.Rows
						}
						continue
					}
				}
			}

			if len(msg) > 0 {
				inputBuf = append(inputBuf, msg...)
				if len(inputBuf) > maxCommandBuffer {
					inputBuf = inputBuf[:0]
					writeError("command buffer exceeded")
					continue
				}
				for {
					idx := bytes.IndexByte(inputBuf, 0)
					if idx < 0 {
						break
					}
					cmdBytes := inputBuf[:idx]
					inputBuf = inputBuf[idx+1:]
					command := strings.TrimSpace(string(cmdBytes))
					if command == "" {
						continue
					}

					ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
					output, err := runHostCommand(ctx, command)
					cancel()
					if output != "" {
						writeText([]byte(output))
					}
					if err != nil {
						writeError(err.Error())
					}
				}
			}
		}
	})
}

func runHostCommand(ctx context.Context, command string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", nil
	}

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx, "cmd.exe", "/C", command)
	} else {
		c = exec.CommandContext(ctx, "/bin/sh", "-c", command)
	}
	out, err := c.CombinedOutput()
	return string(out), err
}
