package websocket

// Hub maintains the set of active websocket clients
// broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub creates new instance of hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// run runs Hub
// unexported method as it is run as soon as websocket handler
// is created
func (h *Hub) run() {
	for {
		select {
		case cli := <-h.register:
			h.clients[cli] = true
		case cli := <-h.unregister:
			if _, ok := h.clients[cli]; ok {
				delete(h.clients, cli)
				close(cli.send)
			}
		case msg := <-h.broadcast:
			for cli := range h.clients {
				select {
				case cli.send <- msg:
				default:
					close(cli.send)
					delete(h.clients, cli)
				}
			}
		}
	}
}
