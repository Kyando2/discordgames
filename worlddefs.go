package main

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
