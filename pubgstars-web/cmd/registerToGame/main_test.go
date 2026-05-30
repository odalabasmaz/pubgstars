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

func TestRegisterToGame_Success(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn:     func(string) model.User { return model.User{Id: "u1", Balance: 100} },
		GetGameByIdFn:        func(string) model.Game { return model.Game{Id: "g1", Price: 50} },
		RegisterUserToGameFn: func(model.User, model.Game) error { return nil },
	}}

	resp, err := app.Handle(context.Background(), event(map[string]interface{}{"id": "g1"}))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d: %s", resp.StatusCode, resp.ErrorMessage)
	}
}

func TestRegisterToGame_StoreError_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn:     func(string) model.User { return model.User{Id: "u1"} },
		GetGameByIdFn:        func(string) model.Game { return model.Game{Id: "g1"} },
		RegisterUserToGameFn: func(model.User, model.Game) error { return errors.New("no spots") },
	}}

	resp, _ := app.Handle(context.Background(), event(map[string]interface{}{"id": "g1"}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	if resp.ErrorMessage != "no spots" {
		t.Errorf("unexpected error message: %q", resp.ErrorMessage)
	}
}
