package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

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
