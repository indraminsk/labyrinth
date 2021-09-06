package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	MaxDimension = 100

	MinValue = 0
	MaxValue = 4
)

type LevelType struct {
	DB      *sql.DB `json:"-"`
	TX      *sql.Tx `json:"-"`
	Creator string
	Game    string
	Level   int64
	Data    [][]int
}

func (obj *LevelType) Validate() (status error) {
	var (
		lenLine int
	)

	// check single line length
	if len(obj.Data) > MaxDimension {
		return errors.New(fmt.Sprintf("max count of lines cannot be more than %d", MaxDimension))
	}

	// init level's length by length of first line
	lenLine = len(obj.Data[0])

	for row, line := range obj.Data {
		// check single line length
		if lenLine > MaxDimension {
			return errors.New(fmt.Sprintf("max line's length cannot be more than %d, broken line %d", MaxDimension, row+1))
		}

		// if length current line does not equal to previous line length, than validation failed
		if lenLine != len(line) {
			return errors.New(fmt.Sprintf("level must be rectangular, broken line is %d", row+1))
		}

		for column, value := range line {
			if (value < MinValue) || (value > MaxValue) {
				return errors.New(fmt.Sprintf("level must contains only [0..4] values, broken value %d in point [%d,%d]", value, row+1, column+1))
			}
		}

		lenLine = len(line)
	}

	return status
}

func (obj *LevelType) Store() (levelId int64) {
	var (
		err error

		tx *sql.Tx

		creatorId, gameId int64
		data              []byte
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

	if !obj.dropLevels(gameId, obj.Level) {
		fmt.Println("[error] can't drop levels")
		return -1
	}

	data, err = json.Marshal(obj.Data)
	if err != nil {
		fmt.Println("[error] can't convert level's data to json")
		return -1
	}

	levelId = obj.addLevels(gameId, obj.Level, data)
	if levelId < 1 {
		fmt.Println("[error] can't add level")
		return -1
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("[error] store commit transaction:", err)
		return -1
	}

	return levelId
}

func (obj *LevelType) addCreator(creator string) (id int64) {
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

func (obj *LevelType) getCreatorId(creator string) (id int64) {
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

func (obj *LevelType) addGame(creatorId int64, game string) (id int64) {
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

func (obj *LevelType) getGameId(creatorId int64, game string) (id int64) {
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

func (obj *LevelType) dropLevels(gameId, level int64) bool {
	var (
		err error

		stmt *sql.Stmt
	)

	stmt, err = obj.TX.Prepare("DELETE FROM levels WHERE (game_id = $1) and (level = $2)")
	if err != nil {
		fmt.Println("[error] drop levels prepare:", err)
		return false
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] drop levels clear stmt memory:", err)
		}
	}()

	_, err = stmt.Exec(gameId, level)
	if err != nil {
		fmt.Println("[error] drop levels execute:", err)
		return false
	}

	return true
}

func (obj *LevelType) addLevels(gameId, level int64, data []byte) (id int64) {
	var (
		err error

		stmt *sql.Stmt
	)

	stmt, err = obj.TX.Prepare("INSERT INTO levels (game_id, level, data) VALUES ($1, $2, $3)")
	if err != nil {
		fmt.Println("[error] add levels prepare:", err)
		return -1
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] add levels clear stmt memory:", err)
		}
	}()

	_, err = stmt.Exec(gameId, level, data)
	if err != nil {
		fmt.Println("[error] add levels execute:", err)
		return -1
	}

	return obj.getLevelId(gameId, level)
}

func (obj *LevelType) getLevelId(gameId, level int64) (id int64) {
	var (
		err error

		stmt *sql.Stmt
		row  *sql.Rows
	)

	stmt, err = obj.TX.Prepare("SELECT id FROM levels WHERE (game_id = $1) and (level = $2)")
	if err != nil {
		fmt.Println("[error] get level id prepare:", err)
		return -1
	}

	defer func() {
		if err = stmt.Close(); err != nil {
			fmt.Println("[error] get level id clear stmt memory:", err)
		}
	}()

	row, err = stmt.Query(gameId, level)
	if err != nil {
		fmt.Println("[error] get level id query statement:", err)
		return -1
	}

	defer func() {
		if err = row.Close(); err != nil {
			fmt.Println("[error] get level id clear row memory:", err)
		}
	}()

	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			fmt.Println("[error] get level id scan row:", err)
			return -1
		}

		return id
	}

	return 0
}
