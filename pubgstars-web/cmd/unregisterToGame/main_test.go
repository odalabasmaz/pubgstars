package main

import (
	"context"
	"errors"
	"testing"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/testutil"
)

func event(body map[string]interface{}) svc.RequestEvent {
	return svc.RequestEvent{
		Body:    body,
		Params:  map[string]map[string]string{"header": {"Authorization": ""}},
		Context: map[string]string{},
	}
}

func TestUnregisterToGame_Success(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn:      func(string) model.User { return model.User{Id: "u1"} },
		GetGameByIdFn:         func(string) model.Game { return model.Game{Id: "g1"} },
		UnregisterUserToGameFn: func(model.User, model.Game) error { return nil },
	}}

	resp, err := app.Handle(context.Background(), event(map[string]interface{}{"id": "g1"}))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestUnregisterToGame_GamePassed_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn:      func(string) model.User { return model.User{Id: "u1"} },
		GetGameByIdFn:         func(string) model.Game { return model.Game{Id: "g1"} },
		UnregisterUserToGameFn: func(model.User, model.Game) error {
			return errors.New("cannot cancel this game as the game date has passed")
		},
	}}

	resp, _ := app.Handle(context.Background(), event(map[string]interface{}{"id": "g1"}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
