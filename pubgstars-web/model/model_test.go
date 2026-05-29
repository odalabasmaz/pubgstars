package model

import "testing"

func TestUser_GetAvailableBalance(t *testing.T) {
	tests := []struct {
		name    string
		balance float64
		bonus   float64
		want    float64
	}{
		{"balance only", 100, 0, 100},
		{"bonus only", 0, 50, 50},
		{"both", 100, 50, 150},
		{"zero", 0, 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := User{Balance: tc.balance, Bonus: tc.bonus}
			if got := u.GetAvailableBalance(); got != tc.want {
				t.Errorf("got %.2f, want %.2f", got, tc.want)
			}
		})
	}
}

func TestGame_Detail(t *testing.T) {
	g := Game{
		League:   "gold",
		Map:      "erangel",
		Type:     "solo",
		Platform: "pc",
		GameDate: "202501011200",
	}
	want := "gold/erangel/solo/pc @202501011200"
	if got := g.Detail(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
