package frontend

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

var srv *http.Server

func Start(configs *sync.Map) {
	engine := initEngine(configs)
	srv = &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}
	go listen()

}
func listen() {
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("listen", "error", err.Error())
		os.Exit(1)
	}
}
func Shutdown() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return srv.Shutdown(ctx)

}
