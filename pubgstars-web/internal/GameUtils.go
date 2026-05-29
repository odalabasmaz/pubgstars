package internal

import (
	"log"
	"time"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

var location, _ = time.LoadLocation("Europe/Istanbul")

func parseGameDate(game model.Game) (time.Time, bool) {
	t, err := time.ParseInLocation("200601021504", game.GameDate, location)
	if err != nil {
		log.Printf("failed to parse game date %q for game %s: %v", game.GameDate, game.Id, err)
		return time.Time{}, false
	}
	return t, true
}

func IsGameInLastHour(game model.Game) bool {
	gameDate, ok := parseGameDate(game)
	if !ok {
		return false
	}
	currentTime := time.Now().In(location)
	return currentTime.After(gameDate.Add(-1*time.Hour)) && currentTime.Before(gameDate.Add(time.Hour))
}

func IsGameCancellable(game model.Game) bool {
	gameDate, ok := parseGameDate(game)
	if !ok {
		return false
	}
	return time.Now().In(location).Before(gameDate.Add(-1 * time.Hour))
}

func IsGameDateValid(game model.Game) bool {
	gameDate, ok := parseGameDate(game)
	if !ok {
		return false
	}
	return time.Now().In(location).Before(gameDate)
}
