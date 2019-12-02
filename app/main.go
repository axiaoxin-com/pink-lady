package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pink-lady/app/apis"
	"pink-lady/app/router"

	"github.com/spf13/viper"
)

func main() {
	log.Println("[INFO] ============ pink-lady ============")
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("[FATAL] ", err)
	}
	app := router.SetupRouter(workdir, "config")
	apis.RegisterRoutes(app)

	bind := viper.GetString("server.bind")
	srv := &http.Server{
		Addr:           bind,
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] listen: %s\n", err)
		}
	}()
	log.Printf("[INFO] Server is listening and serving HTTP on %s (pid: %d)\n", bind, syscall.Getpid())

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[FATAL] Server Shutdown: ", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("[INFO] timeout of 5 seconds.")
	}
	log.Println("[INFO] Server exiting")
}
