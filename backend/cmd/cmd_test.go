package cmd

import (
	"errors"
	"testing"

	"github.com/matobet/verdi/env"
)

type MyParams struct {
	Value int
}

var MyCustomError = errors.New("Non zero parameter!")

func (params *MyParams) Validate() error {
	if params.Value == 0 {
		return MyCustomError
	}
	return nil
}

func TestCommandHandlerInvokesValidate(t *testing.T) {
	handler := func(backend env.Backend, params *MyParams) {}
	_, err := invokeCommandHandler(nil, handler, map[string]interface{}{
		"Value": 0,
	})

	if err != MyCustomError {
		t.Errorf("Expected parameter validation error: %s", err)
	}
}

type MockBackend struct {
	env.RedisPool
	env.Commander
}

func TestCommandHandlerInvokesHandlerWithParsedParams(t *testing.T) {
	called := false
	handler := func(backend env.Backend, params *MyParams) (interface{}, error) {
		called = true
		if params.Value != 42 {
			t.Errorf("Invalid parameter value parsed: %s", params.Value)
		}
		return "ok", nil
	}
	result, err := invokeCommandHandler(MockBackend{nil, nil}, handler, map[string]interface{}{
		"Value": 42,
	})

	if !called {
		t.Error("Expected command handler to be called")
	}
	if err != nil {
		t.Errorf("Command invocation should not fail: %s", err)
	}
	if result != "ok" {
		t.Errorf("Command invocation should return correct value. Instead got: %s", result)
	}
}
