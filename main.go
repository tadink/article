package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"GinTest/db"
	"GinTest/frontend"
	"GinTest/log"
)

func main() {
	log.InitLog()
	err := db.Init()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	siteConfigs, err := db.LoadConfigs()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	frontend.Start(siteConfigs)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")
	if err := frontend.Shutdown(); err != nil {
		slog.Error("Server Shutdown:", "error", err.Error())
		os.Exit(1)
	}
	slog.Info("Server exiting")

}
