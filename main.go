package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leonard475192/go-stations/db"
	"github.com/leonard475192/go-stations/handler"
	"github.com/leonard475192/go-stations/service"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()
	svc := service.NewTODOService(todoDB)

	// set http handlers
	mux := http.NewServeMux()

	// TODO: ここから実装を行う
	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/todos", handler.NewTODOHandler(svc))

	return http.ListenAndServe(port, mux)
}
