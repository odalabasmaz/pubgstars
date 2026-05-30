package main

import (
	"context"
	"errors"
	"testing"

	svc "github.com/odalabasmaz/pubgstars/pubgstars-web/internal"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/model"
	"github.com/odalabasmaz/pubgstars/pubgstars-web/testutil"
)

func makeEvent(body map[string]interface{}) svc.RequestEvent {
	return svc.RequestEvent{
		Body:    body,
		Params:  map[string]map[string]string{"header": {"Authorization": ""}},
		Context: map[string]string{},
	}
}

func validBody() map[string]interface{} {
	return map[string]interface{}{
		"amount":         "200",
		"iban":           "TR123456",
		"nameSurname":    "John Doe",
		"secretQuestion": "pet",
		"secretAnswer":   "dog",
	}
}

func TestWithdraw_InvalidAmount_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	body := validBody()
	body["amount"] = "notanumber"
	resp, _ := app.Handle(context.Background(), makeEvent(body))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWithdraw_BelowMinimum_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{}}
	body := validBody()
	body["amount"] = "50"
	resp, _ := app.Handle(context.Background(), makeEvent(body))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWithdraw_WrongSecretAnswer_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn: func(string) model.User {
			return model.User{Id: "u1", SecretQuestion: "pet", SecretAnswer: "cat", Balance: 500}
		},
	}}
	body := validBody() // secretAnswer = "dog", user has "cat"
	resp, _ := app.Handle(context.Background(), makeEvent(body))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	if resp.ErrorMessage != "Secret question or answer is incorrect." {
		t.Errorf("unexpected message: %q", resp.ErrorMessage)
	}
}

func TestWithdraw_InsufficientBalance_Returns400(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn: func(string) model.User {
			return model.User{Id: "u1", SecretQuestion: "pet", SecretAnswer: "dog", Balance: 50}
		},
	}}
	resp, _ := app.Handle(context.Background(), makeEvent(validBody()))
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	if resp.ErrorMessage != "Insufficient balance!" {
		t.Errorf("unexpected message: %q", resp.ErrorMessage)
	}
}

func TestWithdraw_Success(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn: func(string) model.User {
			return model.User{Id: "u1", SecretQuestion: "pet", SecretAnswer: "dog", Balance: 500}
		},
		UpdateUserWithTxFn: func(model.User, model.TransactionLog) error { return nil },
	}}
	resp, _ := app.Handle(context.Background(), makeEvent(validBody()))
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d: %s", resp.StatusCode, resp.ErrorMessage)
	}
}

func TestWithdraw_DBError_Returns500(t *testing.T) {
	app := &App{store: &testutil.MockStore{
		GetUserByEmailFn: func(string) model.User {
			return model.User{Id: "u1", SecretQuestion: "pet", SecretAnswer: "dog", Balance: 500}
		},
		UpdateUserWithTxFn: func(model.User, model.TransactionLog) error {
			return errors.New("dynamodb error")
		},
	}}
	resp, _ := app.Handle(context.Background(), makeEvent(validBody()))
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
