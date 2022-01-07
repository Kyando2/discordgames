package discordgames

import (
	_ "gopkg.in/yaml.v2"
)

type character uint8
type empire uint8

type position struct {
	x int32
	y int32
}

type EffectType int

type Compendium struct {
	data []*struct {
		Name    string `yaml:"name"`
		Id      int    `yaml:"id"`
		Edible  bool   `yaml:"edible"`
		Fishing int    `yaml:"fishing"`
		Exploit int    `yaml:"exploit"`
		Effects []*struct {
			Type  EffectType `yaml:"type"`
			Value int        `yaml:"value"`
		} `yaml:"effects"`
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

type DataModel struct {
	userdata map[string]*userData
	worldmap map[position]*tileType
}
