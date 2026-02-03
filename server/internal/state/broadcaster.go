package state

import "sync"

type SSEClient chan []byte

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

func (b *Broadcaster) Send(msg []byte) bool {
	select {
	case b.send <- msg:
		return true
	default:
		return false
	}
}

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
					delete(b.clients, c)
					close(c)
				}
			}
			b.mu.Unlock()
		}
	}
}
