package main

import (
	Model "../model"
	"encoding/json"
	"fmt"
)

func main() {
	game := Model.Game{}
	game.Id = "1"
	res, _ := json.Marshal(game)
	fmt.Println(string(res))

	var dat map[string]interface{}
	json.Unmarshal(res, &dat)
	fmt.Println(dat)

	dat["Registered"] = true
	fmt.Println(dat)
}
