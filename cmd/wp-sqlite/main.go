package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/ronaudinho/wp/internal/handler"
	"github.com/ronaudinho/wp/internal/handler/rest"
	"github.com/ronaudinho/wp/internal/handler/websocket"
	"github.com/ronaudinho/wp/internal/repository/sqlite"
	"github.com/ronaudinho/wp/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./message.db")
	if err != nil {
		panic(err)
	}
	err = sqlite.InitDB(db)
	if err != nil {
		panic(err)
	}
	sqrepo := sqlite.New(db)
	svc := service.New(sqrepo)
	rst := rest.New(svc)
	wsc := websocket.New()
	hndlr := handler.New([]handler.IHandler{rst, wsc})
	rou := mux.NewRouter()
	hndlr.WithRoutes(rou)

	srv := &http.Server{
		Handler: rou,
		Addr:    "127.0.0.1:3195",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server started on :3195")
	<-done
	log.Println("server shutdown commencing")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		db.Close()
		cancel()
	}()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Println("server gracefully shutdown")
}
