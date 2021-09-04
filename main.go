package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
)

type ServerType struct {
	db *sql.DB
}

func (server *ServerType) handler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	_, err = w.Write([]byte(fmt.Sprintf("%d", 200)))
	if err != nil {
		http.Error(w, "internal error", 500)
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

	port = flag.Int("p", 9080, "service port")
	flag.Parse()

	fmt.Println("service run on port", *port)
	fmt.Println("to stop the service, press [Ctrl+C]")

	fmt.Println("config build")

	cfg = conf()
	if cfg == nil {
		return
	}

	fmt.Println("config build: done")
	fmt.Println("connect to db")

	db = openDB(cfg.Database)
	if db == nil {
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("error [clear memory db]", err)
		}
	}()

	server = ServerType{db: db}

	fmt.Println("connect to db: done")

	http.HandleFunc("/", server.handler)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
