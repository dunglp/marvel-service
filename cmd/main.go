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
	os.Setenv("APP_HOST_NAME", "http://gateway.marvel.com")
	os.Setenv("APP_PUBLIC_KEY", "5d04d11a2b5d49b6aa05a49998b5f083")
	os.Setenv("APP_PRIVATE_KEY", "481732840d168b3244a9f2a4efb3c21d2b1ce2e0")
	os.Setenv("APP_TS", "1")
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
