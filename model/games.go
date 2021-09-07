package model

import (
	"database/sql"
	"fmt"
)

// structure describe game
type GameType struct {
	TX        *sql.Tx `json:"-"`
	CreatorId int64
	Game      string
}

// add new game (if it needs). prepare sql statement and execute it.
// return id for specific creator's game
func (obj *GameType) addGame() (id int64) {
	var (
		err error

		stmt *sql.Stmt
	)

	id = obj.getGameId()
	if id < 1 {
		stmt, err = obj.TX.Prepare("INSERT INTO games (creator_id, game) VALUES ($1, $2)")
		if err != nil {
			fmt.Println("[error] add new game prepare:", err)
			return -1
		}

		defer func() {
			if err = stmt.Close(); err != nil {
				fmt.Println("[error] add new game clear stmt memory:", err)
			}
		}()

		_, err = stmt.Exec(obj.CreatorId, obj.Game)
		if err != nil {
			fmt.Println("[error] add new game execute:", err)
			return -1
		}

		id = obj.getGameId()
	}

	return id
}

// prepare sql statement and execute it.
// return id for specific creator's game
func (obj *GameType) getGameId() (id int64) {
	var (
		err error

		stmt *sql.Stmt
		row  *sql.Rows
	)

	stmt, err = obj.TX.Prepare("SELECT id FROM games WHERE (creator_id = $1) and (game = $2)")
	if err != nil {
		fmt.Println("[error] get game id prepare:", err)
		return -1
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] get game id clear stmt memory:", err)
		}
	}()

	row, err = stmt.Query(obj.CreatorId, obj.Game)
	if err != nil {
		fmt.Println("[error] get game id query statement:", err)
		return -1
	}

	defer func() {
		if err = row.Close(); err != nil {
			fmt.Println("[error] get game id clear row memory:", err)
		}
	}()

	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			fmt.Println("[error] get game id scan row:", err)
			return -1
		}

		return id
	}

	return 0
}
