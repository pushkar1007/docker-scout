package api

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"io"
	
	"github.com/gorilla/websocket"
)

type terminalRequest struct {
	Type    string `json:"type"`    // "start" or "exitTerm"
	Command string `json:"command"` // command to execute
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
			log.Printf("terminal ws upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		var (
			mu          sync.Mutex
			currentCmd  *exec.Cmd
			cancelFunc  context.CancelFunc
		)

		// Handle incoming messages from frontend
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				if cancelFunc != nil {
					cancelFunc()
				}
				return
			}

			var req terminalRequest
			if err := json.Unmarshal(message, &req); err != nil {
				sendError(conn, &mu, "invalid request format")
				continue
			}

			switch req.Type {
			case "exitTerm":
				mu.Lock()
				if cancelFunc != nil {
					cancelFunc()
					cancelFunc = nil
				}
				if currentCmd != nil && currentCmd.Process != nil {
					currentCmd.Process.Kill()
					currentCmd = nil
				}
				mu.Unlock()
				sendMessage(conn, &mu, map[string]interface{}{
					"type": "terminated",
					"message": "command terminated",
				})

			case "start":
				if req.Command == "" {
					sendError(conn, &mu, "command is required")
					continue
				}

				// Cancel any existing command
				mu.Lock()
				if cancelFunc != nil {
					cancelFunc()
				}
				if currentCmd != nil && currentCmd.Process != nil {
					currentCmd.Process.Kill()
				}
				mu.Unlock()

				// Execute the command
				ctx, cancel := context.WithCancel(context.Background())
				mu.Lock()
				cancelFunc = cancel
				currentCmd = exec.CommandContext(ctx, "sh", "-c", req.Command)
				mu.Unlock()

				// Get stdout and stderr pipes
				stdout, err := currentCmd.StdoutPipe()
				if err != nil {
					sendError(conn, &mu, "failed to get stdout: "+err.Error())
					cancel()
					continue
				}

				stderr, err := currentCmd.StderrPipe()
				if err != nil {
					sendError(conn, &mu, "failed to get stderr: "+err.Error())
					cancel()
					continue
				}

				// Start the command
				if err := currentCmd.Start(); err != nil {
					sendError(conn, &mu, "failed to start command: "+err.Error())
					cancel()
					continue
				}

				sendMessage(conn, &mu, map[string]interface{}{
					"type": "started",
					"message": "command started",
				})

				// Stream stdout
				go func() {
				    scanner := bufio.NewScanner(stdout)
				    for scanner.Scan() {
				        line := scanner.Text()
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "stdout",
				            "data": line,
				        })
				    }
				}()
				
				// Stream stderr
				go func() {
				    scanner := bufio.NewScanner(stderr)
				    for scanner.Scan() {
				        line := scanner.Text()
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "stderr",
				            "data": line,
				        })
				    }
				}()
				
				// Wait for command completion
				// Buffer and send all at once after completion
				go func() {
				    // Read all available data from pipes into memory
				    outData, _ := io.ReadAll(stdout)
				    errData, _ := io.ReadAll(stderr)
				
				    // Wait for command to finish
				    err := currentCmd.Wait()
				
				    mu.Lock()
				    currentCmd = nil
				    cancelFunc = nil
				    mu.Unlock()
				
				    // Send the accumulated stdout once
				    if len(outData) > 0 {
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "stdout",
				            "data": string(outData),
				        })
				    }
				
				    // Send the accumulated stderr once
				    if len(errData) > 0 {
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "stderr",
				            "data": string(errData),
				        })
				    }
				
				    // Send final exit message
				    if err != nil {
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "exit",
				            "error": err.Error(),
				            "code": currentCmd.ProcessState.ExitCode(),
				        })
				    } else {
				        sendMessage(conn, &mu, map[string]interface{}{
				            "type": "exit",
				            "code": 0,
				            "message": "command completed successfully",
				        })
				    }
				}()

			default:
				sendError(conn, &mu, "unknown request type: "+req.Type)
			}
		}
	})
}

func sendMessage(conn *websocket.Conn, mu *sync.Mutex, data interface{}) {
	mu.Lock()
	defer mu.Unlock()

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		log.Printf("failed to write message: %v", err)
	}
}

func sendError(conn *websocket.Conn, mu *sync.Mutex, message string) {
	sendMessage(conn, mu, map[string]interface{}{
		"type": "error",
		"message": message,
	})
}
