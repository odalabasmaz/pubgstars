package main

import (
	"context"
	"errors"
	"testing"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/testutil"
)

func makeEvent(method string, body map[string]interface{}) svc.RequestEvent {
	return svc.RequestEvent{
		Body:    body,
		Params:  map[string]map[string]string{"header": {"Authorization": ""}},
		Context: map[string]string{"http-method": method},
	}
}

func TestAddBudget_WrongMethod_Returns405(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	resp, _ := app.Handle(context.Background(), makeEvent("GET", nil))
	if resp.StatusCode != 405 {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestAddBudget_InvalidBalance_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	resp, _ := app.Handle(context.Background(), makeEvent("POST", map[string]interface{}{
		"userId": "u1", "balance": "bad", "bonus": "10",
	}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAddBudget_InvalidBonus_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	resp, _ := app.Handle(context.Background(), makeEvent("POST", map[string]interface{}{
		"userId": "u1", "balance": "100", "bonus": "bad",
	}))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAddBudget_Success(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByIdFn:      func(string) model.User { return model.User{Id: "u1", Balance: 0, Bonus: 0} },
		UpdateUserWithTxFn: func(model.User, model.TransactionLog) error { return nil },
	}}
	resp, _ := app.Handle(context.Background(), makeEvent("POST", map[string]interface{}{
		"userId": "u1", "balance": "100", "bonus": "50",
	}))
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d: %s", resp.StatusCode, resp.ErrorMessage)
	}
	user, ok := resp.Body.(model.User)
	if !ok {
		t.Fatal("expected User in body")
	}
	if user.Balance != 100 || user.Bonus != 50 {
		t.Errorf("unexpected balance/bonus: %.2f / %.2f", user.Balance, user.Bonus)
	}
}

func TestAddBudget_DBError_Returns500(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByIdFn:      func(string) model.User { return model.User{Id: "u1"} },
		UpdateUserWithTxFn: func(model.User, model.TransactionLog) error { return errors.New("db error") },
	}}
	resp, _ := app.Handle(context.Background(), makeEvent("POST", map[string]interface{}{
		"userId": "u1", "balance": "100", "bonus": "0",
	}))
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
