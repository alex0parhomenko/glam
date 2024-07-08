package main

import (
	"context"
	"flag"
	"glam/internal/config"
	"glam/internal/db"
	"glam/internal/server"
	"glam/pkg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Version = ""

func main() {
	configPath := flag.String("config", os.Getenv("CONFIG_PATH"), "Path to service config")
	flag.Parse()
	conf := pkg.GetConfig[config.Config](configPath)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	client := db.GetDbClient()

	s := server.SpawnServer(conf, client)
	log.Print("Web server started")
	<-signalCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Print("Web server stopped")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	err := client.Disconnect(ctx)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("Db connections closed")
	}
}
