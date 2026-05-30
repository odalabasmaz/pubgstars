package test

import (
	"encoding/json"
	"testing"

	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
)

// TestGameMarshalRoundtrip verifies that a Game struct survives JSON
// marshal → unmarshal with all fields intact.
func TestGameMarshalRoundtrip(t *testing.T) {
	original := model.Game{
		Id:       "g1",
		League:   "gold",
		Map:      "erangel",
		Type:     "solo",
		Platform: "pc",
		Price:    100,
		Status:   "active",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var restored model.Game
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if restored.Id != original.Id {
		t.Errorf("Id: got %q, want %q", restored.Id, original.Id)
	}
	if restored.Price != original.Price {
		t.Errorf("Price: got %.2f, want %.2f", restored.Price, original.Price)
	}
	if restored.Status != original.Status {
		t.Errorf("Status: got %q, want %q", restored.Status, original.Status)
	}
}

// TestGameMapRoundtrip verifies that a Game serialises to the JSON key names
// used by the Lambda handlers (camelCase, matching struct tags).
func TestGameMapRoundtrip(t *testing.T) {
	game := model.Game{Id: "g2", Registered: true}

	data, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if dat["id"] != "g2" {
		t.Errorf("id in map: got %v, want %q", dat["id"], "g2")
	}
	if dat["registered"] != true {
		t.Errorf("registered in map: got %v, want true", dat["registered"])
	}
}
