package main

import (
	"database/sql"
	"encoding/json"
	"log"
)

func getUserData() (u map[string]userData, err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	row, err := db.Query("SELECT * FROM userdata")
	u = make(map[string]userData, 0)
	if err != nil {
		return
	}

	for row.Next() {
		var x, y, char, rank, emp int
		var invraw, id string
		err = row.Scan(&x, &y, &char, &emp, &invraw, &rank, &id)

		if err != nil {
			return
		}
		var inv Inventory
		json.Unmarshal([]byte(invraw), &inv)

		u[id] = userData{
			pos:  pos(x, y),
			char: character(char),
			emp:  empire(emp),
			inv:  inv,
			rank: rank,
		}
	}
	if len(u) == 0 {
		// for testing purposes
		id := "331431342438875137"
		u[id] = userData{
			pos:  pos(0, 0),
			char: character(0),
			emp:  empire(0),
			inv: Inventory{
				Data: []Item{{
					ID:       0,
					Quantity: 1,
				}},
			},
			rank: 0,
		}
		if err := saveUserData(id, u[id]); err != nil {
			log.Fatalf("Could not save user data because %s", err)
		}
	}
	return
}

// Delegates error handling to caller
func saveUserData(id string, u userData) (err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	insertPositionStmt := `INSERT INTO userdata(x, y, char, emp, inv, rank, id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(insertPositionStmt); err != nil {
		return
	}
	defer stmt.Close()
	mrshld, err := json.Marshal(u.inv)
	if err != nil {
		return
	}
	_, err = stmt.Exec(u.pos.x, u.pos.y, u.char, u.emp, string(mrshld), u.rank, id)
	if err != nil {
		return err
	}
	return
}
