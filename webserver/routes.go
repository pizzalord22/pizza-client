package webserver

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "workspace/pizza-anguish-client/telnet"
    ws "workspace/pizza-anguish-client/websocket"
)

// global websocket upgrader variable that is used to upgrade http connections to websocket connections
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// serve the index.html file
func HomePage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

// upgrade the http connection to a websocket connection and start a telnet connection
func WsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    telnetConn := telnet.InitTelnet("anguish.org", 2222)
    defer telnetConn.Close()

    // create buffers for the communication between the telnet and websocket threads
    var toTelnet = make(chan []byte, 10000)
    var toClient = make(chan []byte, 10000)

    go telnet.Reader(telnetConn, toClient)
    go telnet.Writer(telnetConn, toTelnet)
    go ws.Reader(conn, toTelnet)
    ws.Writer(conn, toClient)
}
