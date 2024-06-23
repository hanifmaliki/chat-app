package websocket

type Hub struct {
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[*Client]bool
	ClientsByID map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
		ClientsByID: make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			h.ClientsByID[client.ID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientsByID, client.ID)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
