package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type GameType struct {
	DB *sql.DB `json:"-"`
	TX *sql.Tx `json:"-"`
	Creator string
	Game    string
	Levels  [][]int
}

func (obj *GameType) Store() (gameId int64) {
	var (
		err error

		tx *sql.Tx

		creatorId int64
		levels []byte
	)

	tx, err = obj.DB.Begin()
	if err != nil {
		fmt.Println("[error] store begin transaction:", err)
		return -1
	}

	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Println("[error] store commit transaction:", err)
		}
	}()

	obj.TX = tx

	creatorId = obj.addCreator(obj.Creator)
	if creatorId < 1 {
		fmt.Println("[error] can't create new creator")
		return -1
	}

	gameId = obj.addGame(creatorId, obj.Game)
	if gameId < 1 {
		fmt.Println("[error] can't create new game")
		return -1
	}

	if !obj.dropLevels(gameId) {
		fmt.Println("[error] can't drop levels")
		return -1
	}

	levels, err = json.Marshal(obj.Levels)
	if err != nil {
		fmt.Println("[error] can't convert levels to json")
		return -1
	}

	if !obj.addLevels(gameId, levels) {
		fmt.Println("[error] can't add levels")
		return -1
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("[error] store commit transaction:", err)
		return -1
	}

	return gameId
}

func (obj *GameType) addCreator(creator string) (id int64) {
	var (
		err error

		stmt *sql.Stmt
	)

	id = obj.getCreatorId(creator)
	if id < 1 {
		stmt, err = obj.TX.Prepare("INSERT INTO creators (creator) VALUES ($1)")
		if err != nil {
			fmt.Println("[error] add new creator prepare:", err)
			return -1
		}

		defer func() {
			if err = stmt.Close(); err != nil {
				fmt.Println("[error] add new creator clear stmt memory:", err)
			}
		}()

		_, err = stmt.Exec(creator)
		if err != nil {
			fmt.Println("[error] add new creator:", err)
			return -1
		}

		id = obj.getCreatorId(creator)
	}

	return id
}

func (obj *GameType) getCreatorId(creator string) (id int64) {
	var (
		err error

		stmt *sql.Stmt
		row  *sql.Rows
	)

	stmt, err = obj.TX.Prepare("SELECT id FROM creators WHERE (creator = $1)")
	if err != nil {
		fmt.Println("[error] get creator id prepare:", err)
		return -1
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] get creator id clear stmt memory:", err)
		}
	}()

	row, err = stmt.Query(creator)
	if err != nil {
		fmt.Println("[error] get creator id query statement:", err)
		return -1
	}

	defer func() {
		if err = row.Close(); err != nil {
			fmt.Println("[error] get creator id clear row memory:", err)
		}
	}()

	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			fmt.Println("[error] get creator id scan row:", err)
			return -1
		}

		return id
	}

	return 0
}

func (obj *GameType) addGame(creatorId int64, game string) (id int64) {
	var (
		err error

		stmt *sql.Stmt
	)

	id = obj.getGameId(creatorId, game)
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

		_, err = stmt.Exec(creatorId, game)
		if err != nil {
			fmt.Println("[error] add new game execute:", err)
			return -1
		}

		id = obj.getGameId(creatorId, game)
	}

	return id
}

func (obj *GameType) getGameId(creatorId int64, game string) (id int64) {
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

	row, err = stmt.Query(creatorId, game)
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

func (obj *GameType) dropLevels(gameId int64) bool {
	var (
		err error

		stmt *sql.Stmt
	)

	stmt, err = obj.TX.Prepare("DELETE FROM levels WHERE (game_id = $1)")
	if err != nil {
		fmt.Println("[error] drop levels prepare:", err)
		return false
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] drop levels clear stmt memory:", err)
		}
	}()

	_, err = stmt.Exec(gameId)
	if err != nil {
		fmt.Println("[error] drop levels execute:", err)
		return false
	}

	return true
}

func (obj *GameType) addLevels(gameId int64, levels []byte) bool {
	var (
		err error

		stmt *sql.Stmt
	)

	stmt, err = obj.TX.Prepare("INSERT INTO levels (game_id, data) VALUES ($1, $2)")
	if err != nil {
		fmt.Println("[error] add levels prepare:", err)
		return false
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] add levels clear stmt memory:", err)
		}
	}()

	_, err = stmt.Exec(gameId, levels)
	if err != nil {
		fmt.Println("[error] add levels execute:", err)
		return false
	}

	return true
}