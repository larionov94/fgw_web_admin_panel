package main

import (
	"context"
	"fgw_web_admin_panel/internal/api"
	"fgw_web_admin_panel/pkg"
	"fgw_web_admin_panel/pkg/logg"
	"fgw_web_admin_panel/pkg/msg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	skipNofS = 4 // skipNofS кол-во пропускаемых кадров стека.
)

func main() {
	logger, _ := logg.NewLogger()
	defer logger.Close()

	err := pkg.LoadEnvFile("")
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	server := api.NewServer(os.Getenv("PORT"), mux, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartServer(ctx); err != nil {
			log.Fatalf("%s: %v", msg.ES5002, err)
		}
	}()
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.StopServer(shutdownCtx); err != nil {
		logger.LogE(msg.ES5003, err, skipNofS)
	}
}
