package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fa7C0n/mdm-be/router"
)

func main() {
	Run()
}

func Run() {
	server := router.NewServer()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down Server ....")

		err := gracefulShutdown(server, 25*time.Second)

		if err != nil {
			log.Printf("server stopped: %v", err.Error())
		}

		os.Exit(0)
	}()

	log.Printf("Listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func gracefulShutdown(server *http.Server, maximumTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), maximumTime)
	defer cancel()

	return server.Shutdown(ctx)
}
