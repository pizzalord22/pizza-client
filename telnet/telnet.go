package telnet

import (
    "fmt"
    "github.com/reiver/go-telnet"
    "io"
    "log"
    "strings"
    "time"
)

// initialize a telnet connection
func InitTelnet(server string, port int) *telnet.Conn {
    conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", server, port))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("got a telnet connection")
    return conn
}

// Reader deals with reading and parsing data from the telnet server
func Reader(conn *telnet.Conn, toClient chan []byte) {
    fmt.Println("start reading telnet data")
    var reader = make([]byte, 1)
    //var row = make([]byte, 2024, 2024)
    //var counter = 0

    for {
        _, err := conn.Read(reader)
        if err != nil {
            if err == io.EOF {
                return
            }
            log.Print("error while reading from telnet server:", err)
            return
        }
        toClient <- reader
        time.Sleep(1 * time.Nanosecond)
        continue
        //if endOfMessage(reader[0], row) {
        //    var msg []byte
        //    if string(row) == "Time out." {
        //        msg = []byte("server closed connection , reload the webpage to reconnect")
        //        return
        //    }
        //
        //    // format the msg and check for special cases (login prompt)
        //    msg = formatMessage(row, counter)
        //    msg = checkLogin(row, msg)
        //    msg = append(msg, reader[0])
        //
        //    // send the data the the websocket server
        //    toClient <- msg
        //
        //    // reset the counter and clear the buffer
        //    counter = 0
        //    for k := range row {
        //        row[k] = 0
        //    }
        //} else {
        //    row[counter] = reader[0]
        //    counter++
        //}
    }
}

// Writer writes data it receives from the websocket connection to the telnet server
func Writer(conn *telnet.Conn, write chan []byte) {
    var err error
    for data := range write {
        data = append(data, []byte("\n")...)
        _, err = conn.Write(data)
        if err != nil {
            close(write)
            log.Println(err)
        }
    }
}

// check if the telnet server is done sending data
func endOfMessage(c1 byte, c2 []byte) bool {
    return string(c1) == "\n" ||
        string(c1) == ">" ||
        strings.Contains(string(c2), "What is your name:") ||
        strings.Contains(string(c2), "Password:") ||
        strings.Contains(string(c2), "Throw the other copy out?")
}

// format the data by creating a new buffer and filling it,
// this is nice because the row buffer is big and often contains a lot of 0 bytes that we do not want to send
func formatMessage(row []byte, counter int) []byte {
    var msg = make([]byte, 0, counter)
    for k, v := range row {
        if k == counter {
            break
        }
        msg = append(msg, v)
    }
    return msg
}

// check for login prompts, we add a new line if there is a login prompt
func checkLogin(row, msg []byte) []byte {
    if strings.Contains(string(row), "What is your name:") {
        for k := range row {
            row[k] = 0
        }
        msg = append(msg, []byte("\n")...)
    }
    if strings.Contains(string(row), "Password:") {
        for k := range row {
            row[k] = 0
        }
        msg = append(msg, []byte("\n")...)
    }
    return msg
}
