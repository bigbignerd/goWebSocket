package main

//Hub 维护所有接入的client,以及消息
type Hub struct {
	clients map[string]*Client
	//接受到用户消息
	message chan *Message
	//register new client to clients
	register chan map[string]*Client
	//unregister client to clients
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		message:    make(chan *Message),
		register:   make(chan map[string]*Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			for u, c := range client {
				h.clients[u] = c
			}
		case message := <-h.message:
			//解析message 获取fromUser,toUser,msg
			if targetClient, ok := h.clients[message.ToUser]; ok {
				select {
				case targetClient.send <- message:
				default:
					close(targetClient.send)
					delete(h.clients, message.ToUser)
				}
			}

		}
	}
}
