package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"inter-protocol-auth-server/internal/application"
	"inter-protocol-auth-server/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	isDone := make(chan os.Signal)
	signal.Notify(isDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	mainApp := application.New(ctx)
	go mainApp.Start()

	<-isDone
	cancel()
	mainApp.Shutdown()
}
