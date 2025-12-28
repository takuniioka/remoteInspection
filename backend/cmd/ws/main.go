package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	appConfig "github.com/inspection-tool/backend/internal/config"
	"github.com/inspection-tool/backend/internal/domain"
	"github.com/joho/godotenv"
)

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
	logger *slog.Logger
}

// Room represents an inspection room with connections
type Room struct {
	InspectionID string
	Connections  map[*Connection]bool
	Broadcast    chan interface{}
	Register     chan *Connection
	Unregister   chan *Connection
	mu            sync.RWMutex
}

// Connection represents a WebSocket connection
type Connection struct {
	ID     string
	UserID string
	Role   domain.Role
	Send   chan interface{}
	Room   *Room
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(logger *slog.Logger) *WebSocketHub {
	return &WebSocketHub{
		rooms:  make(map[string]*Room),
		logger: logger,
	}
}

// GetOrCreateRoom gets or creates a room
func (h *WebSocketHub) GetOrCreateRoom(inspectionID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.rooms[inspectionID]; exists {
		return room
	}

	room := &Room{
		InspectionID: inspectionID,
		Connections:  make(map[*Connection]bool),
		Broadcast:    make(chan interface{}, 256),
		Register:     make(chan *Connection),
		Unregister:   make(chan *Connection),
	}

	h.rooms[inspectionID] = room
	go room.run(h.logger)

	return room
}

// Room run loop
func (r *Room) run(logger *slog.Logger) {
	for {
		select {
		case conn := <-r.Register:
			r.mu.Lock()
			r.Connections[conn] = true
			r.mu.Unlock()
			logger.Info("Connection registered",
				slog.String("inspectionId", r.InspectionID),
				slog.String("userId", conn.UserID),
			)

		case conn := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Connections[conn]; ok {
				delete(r.Connections, conn)
				close(conn.Send)
			}
			r.mu.Unlock()
			logger.Info("Connection unregistered",
				slog.String("inspectionId", r.InspectionID),
				slog.String("userId", conn.UserID),
			)

		case message := <-r.Broadcast:
			r.mu.RLock()
			for conn := range r.Connections {
				select {
				case conn.Send <- message:
				default:
					// Channel full, skip
				}
			}
			r.mu.RUnlock()
		}
	}
}

// BroadcastMessage broadcasts a message to the room
func (r *Room) BroadcastMessage(msg interface{}) {
	select {
	case r.Broadcast <- msg:
	default:
		// Channel full, skip
	}
}

// HandleCaptureRequest handles a capture request message
func (h *WebSocketHub) HandleCaptureRequest(room *Room, msg *domain.CaptureRequest) {
	// Send ACK
	ack := map[string]interface{}{
		"type":               "capture.ack",
		"inspectionId":       msg.InspectionID,
		"captureRequestId":   msg.CaptureRequestID,
		"requestedAt":        msg.RequestedAt,
		"requestedBy":        msg.RequestedBy,
		"state":              msg.State,
	}
	room.BroadcastMessage(ack)
}

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Setup logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := appConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create WebSocket hub
	hub := NewWebSocketHub(logger)

	// For MVP, just log that the WebSocket server would start here
	logger.Info("WebSocket server configured",
		slog.String("port", string(rune(cfg.API.Port))+":8081"),
	)

	// TODO: Implement actual WebSocket server with gorilla/websocket
	// For now, this is a placeholder showing the architecture
	_ = hub
	_ = context.Background()

	// Block forever
	select {}
}
