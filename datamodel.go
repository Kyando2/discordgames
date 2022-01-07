package main

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"os"

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
	items []*struct {
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

type item struct {
	id string
}

type inventory struct {
	data []*item
}

type userData struct {
	pos  *position
	char *character
	emp  *empire
	inv  *inventory
	rank int
}

type tileType int

const (
	// tile types
	flatland tileType = 0
	forest   tileType = 1
	mountain tileType = 2
	hills    tileType = 3
	// add more needed tile types
)
const TERRAINTYPEMAX = 4
const DEFWIDTH = 1e4
const DEFHEIGHT = 1e4

type World struct {
	userdata map[string]*userData
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
		// TODO: HANDLE READING
	} else if errors.Is(err, os.ErrNotExist) {
		m, db, _ := initDBWorld()
	} else {
		return
	}
	return
}

func initDBWorld() (m map[position]tileType, db *sql.DB, err error) {
	// Initialize database
	file, err := os.Create("data/wrld.db")
	if err != nil {
		return
	}
	file.Close()
	db, _ = sql.Open("sqlite3", "./sqlite-database.db")
	defer db.Close()
	genWorldTable(db)
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

func genWorldTable(db *sql.DB) (err error) {
	err = nil
	createTblStmt := `CREATE TABLE worldmap (
		"x" integer NOT NULL
		"y" integer NOT NULL
		"type" integer NOT NULL
	);`
	log.Println("Creating table...")
	stmt, err := db.Prepare(createTblStmt)
	if err != nil {
		return
	}
	stmt.Exec()
	return
}

func NewWorld() (World, error) {
	rawitems, err := os.ReadFile("model/items.yaml")
	var comp Compendium
	// Attempt to parse items.yaml
	if err == nil {
		yaml.Unmarshal(rawitems, &comp)
	} else {
		return World{}, err
	}
	worldmap, err := getMap()
}
