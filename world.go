package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func pos(x, y int) position {
	return position{int32(x), int32(y)}
}

func getMap() (m map[position]tileType, err error) {
	m = make(map[position]tileType)
	// Check if database exists
	if _, err = os.Stat("data/wrld.db"); err == nil {
		m, _ = readMap()
	} else if errors.Is(err, os.ErrNotExist) {
		m, _ = initDBWorld()
		carveMap(m)
	}
	return
}

func (w World) At(x, y int) (i int) {
	i = int(w.worldmap[pos(x, y)])
	return
}

// Save the world map into the database at data/wrld.db
func carveMap(m map[position]tileType) (err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	// Insert into the map through multiple batches
	valueStrings := make([]string, 0, len(m))
	valueArgs := make([]interface{}, 0, len(m)*3)
	i := 0
	j := 0
	for k, v := range m {
		if err != nil {
			log.Fatalf("Error carving map: %s", err)
			return
		}
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, k.x, k.y, v)
		if i == 10000 {
			j++
			if j%100 == 0 {
				log.Println("1 million Hot Batch!")
			}
			stm := fmt.Sprintf("INSERT INTO worldmap(x, y, type) VALUES %s;",
				strings.Join(valueStrings, ","))
			_, err = db.Exec(stm, valueArgs...)
			if err != nil {
				log.Fatalf("Error running statement: %s", err)
				return
			}
			i = 0
			valueStrings = make([]string, 0, 100)
			valueArgs = make([]interface{}, 0, 300)
		}
		i++
	}

	return
}

func readMap() (m map[position]tileType, err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	row, err := db.Query("SELECT * FROM worldmap")
	m = make(map[position]tileType)

	if err != nil {
		log.Fatalf("Error querying from DB: %s", err)
		return
	}

	for row.Next() {
		var x, y, ttype int
		err = row.Scan(&x, &y, &ttype)

		if err != nil {
			return
		}

		m[pos(x, y)] = tileType(ttype)
	}
	return
}

func initDBWorld() (m map[position]tileType, err error) {
	// Initialize database
	file, err := os.Create("data/wrld.db")
	if err != nil {
		return
	}
	file.Close()
	genTables()
	m = createRandWorld(DEFHEIGHT, DEFWIDTH)
	return
}

func createRandWorld(width, height int) (m map[position]tileType) {
	m = make(map[position]tileType)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			m[pos(i, j)] = tileType(rand.Intn(TERRAINTYPEMAX))
		}
	}
	return
}

func genTables() (err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	createStmts := make([]string, 2)
	createStmts[0] = `CREATE TABLE worldmap (
		"x" integer NOT NULL,
		"y" integer NOT NULL,
		"type" integer NOT NULL
		);`
	createStmts[1] = `CREATE TABLE userdata (
		"x" integer NOT NULL,
		"y" integer NOT NULL,
		"char" integer NOT NULL,
		"emp" integer NOT NULL,
		"inv" json,
		"rank" integer NOT NULL,
		"id" VARCHAR(255) NOT NULL
		);`

	log.Println("Creating tables...")

	for _, stmt := range createStmts {
		var stm *sql.Stmt
		if stm, err = db.Prepare(stmt); err != nil {
			log.Fatalf("Error preparing table creation statements: %s", err)
			return
		}
		stm.Exec()
	}

	return
}
