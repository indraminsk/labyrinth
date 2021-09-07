// package main run simple http server
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"greenjade/config"
	"greenjade/database"
	"greenjade/handler"
	"net/http"
)

func main() {
	var (
		err error

		port *int
		cfg  *config.ConfType
		db   *sql.DB

		server handler.ServerType
	)

	fmt.Println("config build...")

	cfg = config.BuildConfig(config.DefaultPath)
	if cfg == nil {
		return
	}

	fmt.Println("config build: done")
	fmt.Println("connect to db...")

	db = database.OpenDB(cfg.Database)
	if db == nil {
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("[error] clear memory db", err)
		}
	}()

	server = handler.ServerType{DB: db, Cfg: cfg}

	fmt.Println("connect to db: done")
	port = flag.Int("p", 9080, "service port")
	flag.Parse()

	fmt.Println("service run on port", *port)
	fmt.Println("to stop the service, press [Ctrl+C]")

	http.HandleFunc("/", server.Handler)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
