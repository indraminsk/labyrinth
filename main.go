package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"greenjade/model"
	"net/http"
)

type ServerType struct {
	db *sql.DB
}

func (server *ServerType) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		decoder *json.Decoder
		game    model.GameType

		resource int64
	)

	fmt.Println()

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("I'm ready to POST only"))
		if err != nil {
			fmt.Println("[error] processing wrong request type:", err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		return
	}

	decoder = json.NewDecoder(r.Body)
	err = decoder.Decode(&game)
	if err != nil {
		fmt.Println("[error] decode request params:", err)
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	fmt.Println("user:", game.Creator)
	fmt.Println("game:", game.Game)
	fmt.Println("levels:", game.Data)

	game.DB = server.db

	resource = game.Store()
	if resource <= 0 {
		fmt.Println("[error] storing game data failed")
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	fmt.Println("resource:", resource)

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fmt.Sprintf("%d", resource)))
	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func main() {
	var (
		err error

		port *int
		cfg *Config
		db *sql.DB

		server ServerType
	)

	fmt.Println("config build...")

	cfg = conf()
	if cfg == nil {
		return
	}

	fmt.Println("config build: done")
	fmt.Println("connect to db...")

	db = openDB(cfg.Database)
	if db == nil {
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("[error] clear memory db", err)
		}
	}()

	server = ServerType{db: db}

	fmt.Println("connect to db: done")
	port = flag.Int("p", 9080, "service port")
	flag.Parse()

	fmt.Println("service run on port", *port)
	fmt.Println("to stop the service, press [Ctrl+C]")

	http.HandleFunc("/", server.handler)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
