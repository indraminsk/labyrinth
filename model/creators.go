// package model provide functionality with db entities
package model

import (
	"database/sql"
	"fmt"
)

// structure describe creator
type CreatorType struct {
	TX      *sql.Tx `json:"-"`
	Creator string
}

// add new creator (if it needs). prepare sql statement and execute it.
// return id for specific creator
func (obj *CreatorType) addCreator() (id int64) {
	var (
		err error

		stmt *sql.Stmt
	)

	id = obj.getCreatorId()
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

		_, err = stmt.Exec(obj.Creator)
		if err != nil {
			fmt.Println("[error] add new creator:", err)
			return -1
		}

		id = obj.getCreatorId()
	}

	return id
}

// prepare sql statement and execute it.
// return id for specific creator
func (obj *CreatorType) getCreatorId() (id int64) {
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

	row, err = stmt.Query(obj.Creator)
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
