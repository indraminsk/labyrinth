// package handler release function to processing http requests
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"greenjade/config"
	"greenjade/model"
	"net/http"
)

// base server structure with global objects such as db and config
type ServerType struct {
	DB  *sql.DB
	Cfg *config.ConfType
}

// filtering request type, decoding request body, validate input json and store json data in db.
// in case errors during those stages response with error code and specific message (if it needs).
// build response with id stored level in db.
func (server *ServerType) Handler(w http.ResponseWriter, r *http.Request) {
	var (
		err, status error

		decoder *json.Decoder
		level   model.LevelType

		resource int64
	)

	fmt.Println()

	// we wait only POST request
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

	// convert request body to level structure
	decoder = json.NewDecoder(r.Body)
	err = decoder.Decode(&level)
	if err != nil {
		fmt.Println("[error] decode request params:", err)
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	// pass to level's instance db connection
	level.DB = server.DB

	// before store we need do some validation
	status = level.Validate(server.Cfg.Constraints)
	if status != nil {
		fmt.Println("[error] level is not valid:", status.Error())
		http.Error(w, status.Error(), http.StatusUnprocessableEntity)

		return
	}

	fmt.Println("user:", level.Creator)
	fmt.Println("game:", level.Game)
	fmt.Println("level:", level.Level)
	fmt.Println("data:", level.Data)

	// store level data only if it's correct
	resource = level.Store()
	if resource <= 0 {
		fmt.Println("[error] storing level data failed")
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	fmt.Println("resource:", resource)

	// prepare response
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fmt.Sprintf("%d", resource)))
	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}
