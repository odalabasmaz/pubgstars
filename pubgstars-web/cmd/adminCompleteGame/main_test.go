package main

import (
	"context"
	"errors"
	"testing"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/testutil"
)

func postEvent(body map[string]interface{}) svc.RequestEvent {
	return svc.RequestEvent{
		Body:    body,
		Params:  map[string]map[string]string{"header": {"Authorization": ""}},
		Context: map[string]string{"http-method": "POST"},
	}
}

func TestCompleteGame_WrongMethod_Returns405(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	ev := svc.RequestEvent{Context: map[string]string{"http-method": "GET"}, Params: map[string]map[string]string{"header": {}}}
	resp, _ := app.Handle(context.Background(), ev)
	if resp.StatusCode != 405 {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestCompleteGame_MissingGameId_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	resp, _ := app.Handle(context.Background(), postEvent(map[string]interface{}{
		// gameId intentionally omitted
		"firstWinner": "u1", "secondWinner": "u2", "thirdWinner": "u3",
	}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCompleteGame_MissingWinner_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetGameByIdFn: func(string) model.Game { return model.Game{Id: "g1"} },
	}}
	resp, _ := app.Handle(context.Background(), postEvent(map[string]interface{}{
		"gameId": "g1",
		// firstWinner missing
		"secondWinner": "u2", "thirdWinner": "u3",
	}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCompleteGame_Success(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetGameByIdFn: func(string) model.Game {
			return model.Game{Id: "g1", Award1st: 500, Award2nd: 300, Award3rd: 200}
		},
		GetUserByIdFn:  func(string) model.User { return model.User{Id: "u1"} },
		CompleteGameFn: func(string, model.Game, model.User, model.User, model.User) error { return nil },
	}}
	resp, _ := app.Handle(context.Background(), postEvent(map[string]interface{}{
		"gameId": "g1", "firstWinner": "u1", "secondWinner": "u2", "thirdWinner": "u3",
	}))
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d: %s", resp.StatusCode, resp.ErrorMessage)
	}
}

func TestCompleteGame_StoreError_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetGameByIdFn:  func(string) model.Game { return model.Game{Id: "g1"} },
		GetUserByIdFn:  func(string) model.User { return model.User{} },
		CompleteGameFn: func(string, model.Game, model.User, model.User, model.User) error {
			return errors.New("transaction failed")
		},
	}}
	resp, _ := app.Handle(context.Background(), postEvent(map[string]interface{}{
		"gameId": "g1", "firstWinner": "u1", "secondWinner": "u2", "thirdWinner": "u3",
	}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
