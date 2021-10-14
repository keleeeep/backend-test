/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.42
 */

package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	handlerUser "github.com/keleeeep/test/internal/pkg/handler/user"
	"github.com/keleeeep/test/internal/pkg/resource/db"
	usecaseUser "github.com/keleeeep/test/internal/pkg/usecase/user"
	"github.com/keleeeep/test/internal/pkg/utils/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	persistentDB, err := db.NewPersistent(cfg.Database.Datasource)
	if err != nil {
		log.Fatalf("persistent db: %v", err)
	}

	userUsecase := usecaseUser.NewUsecase(persistentDB,cfg.SecretKey.AccessSecret)
	userHandler := handlerUser.NewHandler(userUsecase)

	r := mux.NewRouter()

	r.HandleFunc("/register", userHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/check-token", userHandler.CheckToken).Methods(http.MethodPost)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		osSignal := <-c

		log.Printf("got %s signal", osSignal.String())

		cancel()
	}()

	if err := startAPI(ctx, cfg.Server.Port, r); err != nil {
		log.Fatalf("start api failed: %v", err)
	}
}

func startAPI(ctx context.Context, addr string, handler http.Handler) error {
	srv := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	log.Printf("server running at %s", addr)

	// wait for context cancellation
	<-ctx.Done()

	log.Printf("shutting down server...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("shutdown server failed: %v", err)
	}

	log.Printf("server stopped gracefully.")
	return nil
}
