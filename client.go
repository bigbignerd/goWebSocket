package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	hub *Hub
	//web socket connection
	conn *websocket.Conn
	//channel缓存发送
	send chan *Message
}

//读消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	//ping message and pong message ???
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message *Message //这里会有潜在的问题吗？？？
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		log.Printf("message: %v", message)
		c.hub.message <- message
	}
}

//写消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//ws server 处理websocket 请求
func serverWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	token := userToken(r)
	if token == "" {
		log.Println("未指定token")
		return
	}
	//create new client
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan *Message, 256), //这里是否应该使用指针?
	}
	//注册new client to hub
	client.hub.register <- map[string]*Client{token: client}
	go client.readPump()
	go client.writePump()
}

//在websockt 握手连接中获取subprotocols，用做用户标识
func userToken(r *http.Request) string {
	subprotocols := websocket.Subprotocols(r)
	//获取用户token
	var token string = ""
	if len := len(subprotocols); len > 1 && subprotocols[0] == "token" {
		token = subprotocols[1]
	}
	return token
}

type User struct {
	UserName string `json:"userName"`
}

//在线用户
func onlineUser(hub *Hub, w http.ResponseWriter, r *http.Request) {
	clients := []User{}
	for u, _ := range hub.clients {
		user := User{
			UserName: u,
		}
		clients = append(clients, user)
	}
	err := json.NewEncoder(w).Encode(clients)

	if err != nil {
		log.Printf("json返回出错：%v", err)
		return
	}
	log.Println(clients)
}
