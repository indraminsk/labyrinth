package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"greenjade/config"
)

// structure describe input data about level and processed data
type LevelType struct {
	DB        *sql.DB `json:"-"`
	TX        *sql.Tx `json:"-"`
	CreatorId int64 `json:"-"`
	GameId    int64 `json:"-"`
	JsonData  []byte `json:"-"`
	Creator   string
	Game      string
	Level     int64
	Data      [][]int
}

// apply to level data constraints. constraints specify in config file section Constraints.
// return nil or error object
func (obj *LevelType) Validate(constraints config.ConstraintsType) (status error) {
	var (
		lenLine int
	)

	// check single line length
	if len(obj.Data) > constraints.Dimension.Max {
		return errors.New(fmt.Sprintf("max count of lines cannot be more than %d", constraints.Dimension.Max))
	}

	// init level's length by length of first line
	lenLine = len(obj.Data[0])

	for row, line := range obj.Data {
		// check single line length
		if lenLine > constraints.Dimension.Max {
			return errors.New(fmt.Sprintf("max line's length cannot be more than %d, broken line %d", constraints.Dimension.Max, row+1))
		}

		// if length current line does not equal to previous line length, than validation failed
		if lenLine != len(line) {
			return errors.New(fmt.Sprintf("level must be rectangular, broken line is %d", row+1))
		}

		// in each column must be only valid integer marks
		for column, value := range line {
			if (value < constraints.Point.Min) || (value > constraints.Point.Max) {
				return errors.New(fmt.Sprintf("level must contains only [0..4] values, broken value %d in point [%d,%d]", value, row+1, column+1))
			}
		}

		// store length of current line to compare with the next line
		lenLine = len(line)
	}

	return status
}

// all needed actions to store level: find or create creator and game, delete previous level data and store new data.
// return id new db's record
func (obj *LevelType) Store() (levelId int64) {
	var (
		err error

		tx *sql.Tx

		creator CreatorType
		game    GameType

		creatorId, gameId int64
	)

	// run transaction common for all storing stages
	tx, err = obj.DB.Begin()
	if err != nil {
		fmt.Println("[error] store begin transaction:", err)
		return -1
	}

	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Println("[error] store rollback transaction:", err)
		}
	}()

	obj.TX = tx

	// init creator structure and create (if it needs) new entity
	creator = CreatorType{TX: tx, Creator: obj.Creator}

	creatorId = creator.addCreator()
	if creatorId < 1 {
		fmt.Println("[error] can't create new creator")
		return -1
	}

	// init game structure and create (if it needs) new entity
	game = GameType{TX: tx, CreatorId: creatorId, Game: obj.Game}

	gameId = game.addGame()
	if gameId < 1 {
		fmt.Println("[error] can't create new game")
		return -1
	}

	// add creator and game id to level structure, drop previous level data
	obj.CreatorId = creatorId
	obj.GameId = gameId

	if !obj.dropLevels() {
		fmt.Println("[error] can't drop levels")
		return -1
	}

	// before storing convert to json level data and add it
	obj.JsonData, err = json.Marshal(obj.Data)
	if err != nil {
		fmt.Println("[error] can't convert level's data to json")
		return -1
	}

	levelId = obj.addLevels()
	if levelId < 1 {
		fmt.Println("[error] can't add level")
		return -1
	}

	// commit common transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("[error] store commit transaction:", err)
		return -1
	}

	return levelId
}

// drop previous level data. prepare sql statement and execute it.
func (obj *LevelType) dropLevels() bool {
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

	_, err = stmt.Exec(obj.GameId, obj.Level)
	if err != nil {
		fmt.Println("[error] drop levels execute:", err)
		return false
	}

	return true
}

// add actual level data. prepare sql statement and execute it.
// return id for record with new data in db
func (obj *LevelType) addLevels() (id int64) {
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

	_, err = stmt.Exec(obj.GameId, obj.Level, obj.JsonData)
	if err != nil {
		fmt.Println("[error] add levels execute:", err)
		return -1
	}

	return obj.getLevelId()
}

// prepare sql statement and execute it.
// return id for specific game and level
func (obj *LevelType) getLevelId() (id int64) {
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

	row, err = stmt.Query(obj.GameId, obj.Level)
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
