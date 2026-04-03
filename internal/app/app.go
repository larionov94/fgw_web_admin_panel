package app

import (
	"context"
	"fgw_web_admin_panel/internal/api"
	"fgw_web_admin_panel/internal/api/middleware"
	"fgw_web_admin_panel/internal/config"
	"fgw_web_admin_panel/internal/config/db"
	"fgw_web_admin_panel/internal/handler/http_web"
	"fgw_web_admin_panel/internal/repository"
	"fgw_web_admin_panel/internal/service"
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

func StartApplication() {
	api.InitSessionStore()

	logger, err := logg.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	authMiddleware := middleware.NewAuthMiddleware(api.Store, logger)

	if err = pkg.LoadEnvFile("", logger); err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.ES5005)
		log.Fatal(err)
	}

	cfgMSSQL, err := config.NewCfgMSSQL(logger)
	if err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB504)
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mssqlDB, err := db.NewConnMSSQL(ctx, cfgMSSQL, logger)
	if err != nil {
		logger.LogEf(logg.SkipNofS, err, "%s", msg.EDB505)
		log.Fatal(err)
	}

	defer db.Close(mssqlDB, logger)
	repoPerformer := repository.NewPerformerRepo(mssqlDB, logger)
	repoRole := repository.NewRoleRepo(mssqlDB, logger)
	repoHistory := repository.NewHistoryRepo(mssqlDB, logger)

	servicePerformer := service.NewPerformerService(repoPerformer, logger)
	serviceRole := service.NewRoleService(repoRole, logger)
	serviceHistory := service.NewHistoryService(repoHistory, logger)

	handlerPerformerAuth := http_web.NewAuthHandler(servicePerformer, serviceHistory, logger, authMiddleware)
	handlerPerformer := http_web.NewPerformerHandler(servicePerformer, serviceRole, logger, authMiddleware)

	mux := http.NewServeMux()

	handlerPerformerAuth.ServeHTTPRouter(mux)
	handlerPerformer.ServeHTTPRouter(mux)

	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))

	server := api.NewServer(os.Getenv("PORT"), mux, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartServer(ctx); err != nil {
			logger.LogEf(logg.SkipNofS, err, "%s", msg.ES5002)
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second)
	logger.LogIf(logg.SkipNofS, "%s: %s", msg.IS2001, os.Getenv("PORT"))

	<-quit

	logger.LogWf(logg.SkipNofS, "%s", msg.WL4002)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.StopServer(shutdownCtx); err != nil {
		logger.LogE(msg.ES5003, err, logg.SkipNofS)
	}
}
