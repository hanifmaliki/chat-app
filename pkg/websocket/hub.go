package websocket

type Hub struct {
	Broadcast         chan []byte
	Register          chan *Client
	Unregister        chan *Client
	Clients           map[*Client]bool
	ClientsByUsername map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:         make(chan []byte),
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Clients:           make(map[*Client]bool),
		ClientsByUsername: make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			h.ClientsByUsername[client.Username] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientsByUsername, client.Username)
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
