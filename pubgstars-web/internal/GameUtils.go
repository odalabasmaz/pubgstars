package internal

import (
	"time"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

var location, _ = time.LoadLocation("Europe/Istanbul")

func IsGameInLastHour(game model.Game) bool {
	currentTime := time.Now().In(location)
	gameDate, _ := time.ParseInLocation("200601021504", game.GameDate, location)
	gameDateMinus1h := gameDate.Add(-1 * time.Hour)
	gameDatePlus1h := gameDate.Add(+1 * time.Hour)
	return currentTime.After(gameDateMinus1h) && currentTime.Before(gameDatePlus1h)
}

func IsGameCancellable(game model.Game) bool {
	currentTime := time.Now().In(location)
	gameDate, _ := time.ParseInLocation("200601021504", game.GameDate, location)
	gameDateMinus1h := gameDate.Add(-1 * time.Hour)
	return currentTime.Before(gameDateMinus1h)
}

func IsGameDateValid(game model.Game) bool {
	currentTime := time.Now().In(location)
	gameDate, _ := time.ParseInLocation("200601021504", game.GameDate, location)
	return currentTime.Before(gameDate)
}
