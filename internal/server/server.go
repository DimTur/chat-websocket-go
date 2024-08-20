package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var server *http.Server

func Run(host, port string, handler http.Handler) error {
	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: 200 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
	}

	return server.ListenAndServe()
}

func Shutdown() error {
	return server.Shutdown(context.Background())
}
