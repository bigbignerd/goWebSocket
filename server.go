package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{Subprotocols: []string{"token"}}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		log.Println(websocket.Subprotocols(r))
		//如果前端close掉了连接 ，server端继续获取数据会产生一个panic：repeated read on failed websocket connection
		//ws.close();所以server端需要做异常处理
		go func(conn *websocket.Conn) {
			for {
				mType, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					return
				}
				conn.WriteMessage(mType, msg)
			}
		}(conn)
	})
	http.ListenAndServe(":3000", nil)
}
