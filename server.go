package main

import (
	"github.com/gorilla/websocket"
	// "log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"token"},
	CheckOrigin:  func(r *http.Request) bool { return true },
}

func main() {

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "t.html")
	})
	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
		// var conn, _ = upgrader.Upgrade(w, r, nil)
		// subprotocols := websocket.Subprotocols(r)
		// //获取用户token
		// var token string = ""
		// if len := len(subprotocols); len > 1 && subprotocols[0] == "token" {
		// 	log.Println(subprotocols[1])
		// 	token = subprotocols[1]
		// } else {
		// 	log.Println("未指定用户token")
		// 	return
		// }
		//如果前端close掉了连接 ，server端继续获取数据会产生一个panic：repeated read on failed websocket connection
		//ws.close();所以server端需要做异常处理
		// go func(conn *websocket.Conn) {
		// 	for {
		// 		mType, msg, err := conn.ReadMessage()
		// 		if err != nil {
		// 			conn.Close()
		// 			return
		// 		}
		// 		conn.WriteMessage(mType, msg)
		// 	}
		// }(conn)
	})
	http.ListenAndServe(":3000", nil)
}
