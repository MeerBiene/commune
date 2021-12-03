package main

import (
	"commune/client"
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

	go func() {
		<-sc

		log.Println("Shutting down server")
		os.Exit(1)
	}()

	config := flag.String("config", "config.toml", "Commune configuration file")

	flag.Parse()

	client.Start(&client.StartRequest{
		Config: *config,
	})
}
