package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leonard475192/go-stations/db"
	"github.com/leonard475192/go-stations/handler"
	"github.com/leonard475192/go-stations/model"
	"github.com/leonard475192/go-stations/service"
)

var Db *sql.DB

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
	Db = todoDB

	// set http handlers
	mux := http.NewServeMux()

	// TODO: ここから実装を行う
	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/todos", http.HandlerFunc(todo))
	log.Fatal(http.ListenAndServe(port, mux))

	return nil
}

func todo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		svc := service.NewTODOService(Db)
		TodoHandler := handler.NewTODOHandler(svc)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		var req model.CreateTODORequest
		json.Unmarshal(body, &req)
		log.Print(req)

		res, err := TodoHandler.Create(ctx, &req)
		if err != nil {
			w.WriteHeader(400)
			log.Print(err)
		}
		res_json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(400)
			log.Print(err)
		}

		w.Write(res_json)
	}
}
