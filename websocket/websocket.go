package ws

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    "log"
    "regexp"
    "strings"
)

// message represents a single message that is send between the websocket server and and the js front end
type message struct {
    DataType string `json:"type"`
    Text     string `json:"text"`
}

// this regex allows us ti identify chat messages,
var chatRegex = regexp.MustCompile("(say|tell|ask|shout|whisper|\\(:\\w*:\\)|<\\w*>|\\[:\\w*:\\]|\\*:\\w*:\\*)")

// The reader function reads incoming websocket messages and sends them to the telnet server when the message type is command
// this is setup this way so it can be easily expanded to allow for more types of messages
func Reader(conn *websocket.Conn, toTelnet chan []byte) {
    for {
        _, m, err := conn.ReadMessage()
        if err != nil {
            if err == websocket.ErrCloseSent || websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway){
                return
            }
            close(toTelnet)
            fmt.Println("error1", err)
        }
        if string(m) == "" {
            return
        }
        var msg message
        err = json.Unmarshal(m, &msg)
        if err != nil {
            fmt.Println("error2", err)
            return
        }
        switch msg.DataType {
        case "command":
            toTelnet <- []byte(msg.Text)
        default:
            fmt.Println("incoming message:", msg.Text)
        }
    }
}

// the writer gets data from the telnet server and sends it to the js frontend
//  this is also where we check if a function is a chat message or a game message
func Writer(conn *websocket.Conn, readChan chan []byte) {
    var err error
    var nextIsChat bool
    var msg message
    for data := range readChan {
        if data[0] != 32 {
            nextIsChat = false
        }
        msg.DataType = "gameText"

        if chatRegex.MatchString(string(data)) || nextIsChat {
            nextIsChat = false
            msg.DataType = "gameChat"
            if strings.Contains(string(data), "\r\n") {
                nextIsChat = true
            }
        }

        msg.Text = string(data)
        err = conn.WriteJSON(msg)
        if err != nil {
            log.Println("ws writer error", err)
            return
        }
    }
}
