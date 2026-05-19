package main

import (
	model "github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"encoding/json"
	"fmt"
)

func main() {
	game := model.Game{}
	game.Id = "1"
	res, _ := json.Marshal(game)
	fmt.Println(string(res))

	var dat map[string]interface{}
	json.Unmarshal(res, &dat)
	fmt.Println(dat)

	dat["Registered"] = true
	fmt.Println(dat)
}
