package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/env"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/go.uuid"
)

type Request struct {
	Name   string                 `json:"name"`
	ID     string                 `json:"id"`
	Params map[string]interface{} `json:"params"`
}

type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

var TimedOut = errors.New("command.Run: Command timed out")

func Run(backend env.Backend, name string, params map[string]interface{}) (result map[string]interface{}, err error) {
	cmd := commands[name]
	if cmd == nil {
		return nil, fmt.Errorf(`command.Run: Unknown command: "%s"`, name)
	}
	request := Request{
		Name:   name,
		ID:     uuid.NewV4().String(),
		Params: params,
	}

	queue, err := queueByClass(cmd.Class, params)
	if err != nil {
		return nil, err
	}

	conn := backend.Redis()
	defer conn.Close()

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, err = conn.Do("LPUSH", queue, requestBody)
	if err != nil {
		return nil, err
	}

	reply := replyQueue(request.ID)
	values, err := redis.Values(conn.Do("BRPOP", reply, config.Conf.CommandTimeout))
	if err != nil || values[0] == nil {
		return nil, TimedOut
	}

	response := values[1].([]byte)

	err = json.Unmarshal(response, &result)
	return
}

func respond(conn env.Redis, requestID string, response *Response) (err error) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		return err
	}
	queue := replyQueue(requestID)

	// TODO: pipeline LPUSH + EXPIRE
	_, err = conn.Do("LPUSH", queue, responseBody)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", queue, config.Conf.CommandTimeout)
	return
}

func respondSuccess(conn env.Redis, requestID string, result interface{}) error {
	return respond(conn, requestID, &Response{
		Status: "success",
		Result: result,
	})
}

func respondError(conn env.Redis, requestID string, result interface{}) error {
	return respond(conn, requestID, &Response{
		Status: "error",
		Result: result,
	})
}

func Listen(backend env.Backend, queue string) {
	conn := backend.Redis()
	defer conn.Close()

	log.Println("Started listening on queue:", queue)

	for {
		values, err := redis.Values(conn.Do("BRPOP", queue, 0))
		if err != nil {
			log.Fatal(err)
		}

		body := values[1].([]byte)

		var request Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			fmt.Println("command.Listen: received malformed response. Skipping.")
			continue
		}

		id := request.ID
		if id == "" {
			fmt.Println("command.Listen: Received a command without 'id'. Skipping")
			continue
		}

		name := request.Name
		if name == "" {
			respondError(conn, id, "Received a command without 'name'")
			continue
		}

		cmd := commands[name]
		if cmd == nil {
			respondError(conn, id, "Unknown command")
			continue
		}

		result, err := invokeCommandHandler(backend, cmd.handler, request.Params)
		if err != nil {
			respondError(conn, id, err.Error())
		} else {
			respondSuccess(conn, id, result)
		}
	}
}

type Validator interface {
	Validate() error
}

func invokeCommandHandler(backend env.Backend, handler interface{}, params map[string]interface{}) (result interface{}, err error) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("cmd: command handler must be a function")
	}
	if handlerType.NumIn() != 2 {
		panic("cmd: command handler must take exactly 2 arguments: env.Backend and params")
	}
	paramPtrType := handlerType.In(1)
	if paramPtrType.Kind() != reflect.Ptr {
		panic("cmd: command handler must take as 2nd argument a pointer to parameter struct")
	}
	cmdParams := reflect.New(paramPtrType.Elem())
	cmdParamsVal := cmdParams.Interface()
	err = mapstructure.Decode(params, cmdParamsVal)
	if err != nil {
		return
	}
	// if params support validation
	if paramValidator, ok := cmdParamsVal.(Validator); ok {
		if err = paramValidator.Validate(); err != nil {
			return
		}
	}
	results := reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(backend), cmdParams})
	result = results[0].Interface()
	err, _ = results[1].Interface().(error)
	return
}
