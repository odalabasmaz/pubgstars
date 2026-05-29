package internal

import (
	"testing"
	"time"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

func gameAt(t time.Time) model.Game {
	return model.Game{GameDate: t.In(location).Format("200601021504")}
}

func TestIsGameDateValid(t *testing.T) {
	tests := []struct {
		name string
		game model.Game
		want bool
	}{
		{"future game is valid", gameAt(time.Now().Add(2 * time.Hour)), true},
		{"past game is invalid", gameAt(time.Now().Add(-2 * time.Hour)), false},
		{"invalid date string", model.Game{GameDate: "bad"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsGameDateValid(tc.game); got != tc.want {
				t.Errorf("IsGameDateValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsGameCancellable(t *testing.T) {
	tests := []struct {
		name string
		game model.Game
		want bool
	}{
		{"far future is cancellable", gameAt(time.Now().Add(3 * time.Hour)), true},
		{"within 1h cutoff is not cancellable", gameAt(time.Now().Add(30 * time.Minute)), false},
		{"past game is not cancellable", gameAt(time.Now().Add(-1 * time.Hour)), false},
		{"invalid date string", model.Game{GameDate: "bad"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsGameCancellable(tc.game); got != tc.want {
				t.Errorf("IsGameCancellable() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsGameInLastHour(t *testing.T) {
	tests := []struct {
		name string
		game model.Game
		want bool
	}{
		{"game starting in 30 min is in window", gameAt(time.Now().Add(30 * time.Minute)), true},
		{"game started 30 min ago is in window", gameAt(time.Now().Add(-30 * time.Minute)), true},
		{"far future is not in window", gameAt(time.Now().Add(3 * time.Hour)), false},
		{"long past is not in window", gameAt(time.Now().Add(-3 * time.Hour)), false},
		{"invalid date string", model.Game{GameDate: "bad"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsGameInLastHour(tc.game); got != tc.want {
				t.Errorf("IsGameInLastHour() = %v, want %v", got, tc.want)
			}
		})
	}
}
