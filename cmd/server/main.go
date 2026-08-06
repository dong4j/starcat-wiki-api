// Package main starcat-wiki-api 入口。
//
// 装配逻辑在可导出的 server 包中，便于 starcat-api 聚合部署复用。
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/starcat-app/starcat-wiki-api/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("[env] no .env file found, using OS environment only")
	}

	svc, err := server.FromEnv()
	if err != nil {
		log.Fatalf("failed to init wiki server: %v", err)
	}
	defer func() { _ = svc.Close() }()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received %v, shutting down...", sig)
		_ = svc.Close()
		os.Exit(0)
	}()

	svc.StartBackground()
	log.Fatal(http.ListenAndServe(svc.Addr(), svc.Handler()))
}
