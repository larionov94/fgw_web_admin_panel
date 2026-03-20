package main

import (
	"context"
	"fgw_web_admin_panel/internal/api"
	"fgw_web_admin_panel/internal/config"
	"fgw_web_admin_panel/internal/config/db"
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
	logger, err := logg.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	if err = pkg.LoadEnvFile("", logger); err != nil {
		logger.LogEf(skipNofS, err, "%s", msg.ES5005)
		log.Fatal(err)
	}

	cfgMSSQL, err := config.NewCfgMSSQL(logger)
	if err != nil {
		logger.LogEf(skipNofS, err, "%s", msg.EDB504)
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mssqlDB, err := db.NewConnMSSQL(ctx, cfgMSSQL, logger)
	if err != nil {
		log.Fatalf("%s: %s", msg.EDB505, err)
	}

	defer db.Close(mssqlDB, logger)

	mux := http.NewServeMux()
	server := api.NewServer(os.Getenv("PORT"), mux, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartServer(ctx); err != nil {
			log.Fatalf("%s: %v", msg.ES5002, err)
		}
	}()

	time.Sleep(time.Second)
	logger.LogIf(skipNofS, "%s: %s", msg.IS2001, os.Getenv("PORT"))

	<-quit

	logger.LogWf(skipNofS, "%s", msg.WL4002)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.StopServer(shutdownCtx); err != nil {
		logger.LogE(msg.ES5003, err, skipNofS)
	}
}
