package state

import "sync"

// SSEClient is a channel for sending JSON events to Server-Sent Events subscribers.
type SSEClient chan []byte

// Broadcaster manages concurrent publish-subscribe for real-time events.
// It serializes Add, Remove, and Send operations through a single goroutine (run)
// to guarantee safe concurrent access without explicit locking on the clients map.
type Broadcaster struct {
	mu      sync.Mutex
	clients map[SSEClient]struct{}
	add     chan SSEClient
	remove  chan SSEClient
	send    chan []byte
}

func NewBroadcaster(buffer int) *Broadcaster {
	b := &Broadcaster{
		clients: make(map[SSEClient]struct{}),
		add:     make(chan SSEClient),
		remove:  make(chan SSEClient),
		send:    make(chan []byte, buffer),
	}
	go b.run()
	return b
}

func (b *Broadcaster) Add(c SSEClient) {
	b.add <- c
}

func (b *Broadcaster) Remove(c SSEClient) {
	b.remove <- c
}

// Send attempts to queue a message for broadcast. Returns false if the send buffer is full,
// indicating the broadcaster is overloaded and clients may drop events.
func (b *Broadcaster) Send(msg []byte) bool {
	select {
	case b.send <- msg:
		return true
	default:
		return false
	}
}

// run serializes all client operations on a single goroutine to avoid concurrent map access.
// Clients that cannot receive (full channel) are immediately dropped to prevent the broadcaster
// from blocking on any slow subscriber.
func (b *Broadcaster) run() {
	for {
		select {
		case c := <-b.add:
			b.mu.Lock()
			b.clients[c] = struct{}{}
			b.mu.Unlock()
		case c := <-b.remove:
			b.mu.Lock()
			delete(b.clients, c)
			close(c)
			b.mu.Unlock()
		case msg := <-b.send:
			b.mu.Lock()
			for c := range b.clients {
				select {
				case c <- msg:
				default:
					// Client channel full; drop to prevent broadcaster stall.
					delete(b.clients, c)
					close(c)
				}
			}
			b.mu.Unlock()
		}
	}
}
