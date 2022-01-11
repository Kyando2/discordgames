package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type character uint8
type empire uint8

type position struct {
	x int32
	y int32
}

type EffectType int

type Compendium struct {
	items []struct {
		Name    string `yaml:"name"`
		Id      int    `yaml:"id"`
		Edible  bool   `yaml:"edible"`
		Fishing int    `yaml:"fishing"`
		Exploit int    `yaml:"exploit"`
		Effects []*struct {
			Type  EffectType `yaml:"type"`
			Value int        `yaml:"value"`
		} `yaml:"effects"`
		Recipes []*struct {
			ResultAmount  int        `yaml:"amount"`
			RequiredItems [][]string `yaml:"required"`
		}
	}
}

type Item struct {
	ID       int
	Quantity int
}

type Inventory struct {
	Data []Item
}

type userData struct {
	pos  position
	char character
	emp  empire
	inv  Inventory
	rank int
}

type tileType int

const (
	// tile types
	flatland tileType = 0
	forest   tileType = 1
	mountain tileType = 2
	hills    tileType = 3
	mine     tileType = 4
	// add more needed tile types
)
const TERRAINTYPEMAX = 10
const DEFWIDTH = 2e3
const DEFHEIGHT = 2e3

type World struct {
	userdata map[string]userData
	worldmap map[position]tileType
	comp     Compendium
}

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
		u["331431342438875137"] = userData{
			pos:  pos(0, 0),
			char: character(0),
			emp:  empire(0),
			inv: Inventory{
				Data: []Item{Item{
					ID:       0,
					Quantity: 1,
				}},
			},
			rank: 0,
		}
	}
	return
}

func (w World) At(x, y int) (i int) {
	i = int(w.worldmap[pos(x, y)])
	return
}

func carveMap(m map[position]tileType) (err error) {
	db, _ := sql.Open("sqlite3", "data/wrld.db")
	insertPositionStmt := `INSERT INTO worldmap(x, y, type) VALUES (?, ?, ?)`
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(insertPositionStmt); err != nil {
		log.Fatalf("Error preparing world map insert statement: %s", err)
		return
	}
	defer stmt.Close()
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

func NewWorld() (w World, err error) {
	rawitems, err := os.ReadFile("model/items.yaml")
	var comp Compendium
	// Attempt to parse items.yaml
	if err == nil {
		if err = yaml.Unmarshal(rawitems, &comp.items); err != nil {
			log.Fatalf("Error unmarshaling yaml: %s", err)
		}
	} else {
		return World{}, err
	}
	var worldmap map[position]tileType
	var userdata map[string]userData
	if worldmap, err = getMap(); err != nil {
		log.Fatalf("Error loading map %s", err)
		return
	}
	if userdata, err = getUserData(); err != nil {
		log.Fatal("Error getting userdata map")
		return
	}
	w = World{
		userdata: userdata,
		worldmap: worldmap,
		comp:     comp,
	}
	return
}
