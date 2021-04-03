package webserver

import (
    "net/http"
    "time"
)

// create a new http server
func InitServer(addr string, port string, timeout time.Duration) *http.Server {
    return &http.Server{
        Addr:         addr + ":" + port,
        IdleTimeout:  timeout,
        ReadTimeout:  timeout,
        WriteTimeout: timeout,
    }
}
