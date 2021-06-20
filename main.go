package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "net/http"
    "os"
    "os/signal"
    "time"
    "workspace/pizza-anguish-client/webserver"
)

// in main we create a basic webserver and start it
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", webserver.HomePage)
    router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "favicon.ico")
    })

    router.HandleFunc("/ws", webserver.WsHandler)
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
    server := webserver.InitServer("", "6565", 30*time.Second)
    server.Handler = router

    go func() {
        err := server.ListenAndServe()
        if err != nil {
            fmt.Println(err)
            return
        }
    }()
    fmt.Println("server started")

    var exit = make(chan os.Signal)
    signal.Notify(exit, os.Interrupt, os.Kill)
    <-exit
}
