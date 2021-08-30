package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leonard475192/go-stations/db"
	"github.com/leonard475192/go-stations/model"
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

	// set http handlers
	mux := http.NewServeMux()

	// TODO: ここから実装を行う
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	log.Fatal(http.ListenAndServe(port, mux))

	return nil
}

func healthz(writer http.ResponseWriter, request *http.Request) {
	msg := model.HealthzResponse{
		Message: "message",
	}
	msg_json, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(msg_json)
}
