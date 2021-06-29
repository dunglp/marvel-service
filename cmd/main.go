package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"xendit-technical-assessment/pkg/service"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Service panicked %s: %s", r, string(debug.Stack()))
			log.Fatalf("Fatal: %v", r)
		}
	}()

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	svc := service.Run()

	go func() {
		<-quit
		svc.Stop()
		done <- true
	}()

	<-done
}
