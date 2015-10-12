package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/util"
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

var redisPool *redis.Pool

func Init() error {
	redisPool = util.NewRedisPool()

	return nil
}

func Run(name string, params map[string]interface{}) (result map[string]interface{}, err error) {
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

	conn := redisPool.Get()
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

func respond(requestID string, response *Response) (err error) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		return err
	}
	queue := replyQueue(requestID)

	conn := redisPool.Get()
	defer conn.Close()

	// TODO: pipeline LPUSH + EXPIRE
	_, err = conn.Do("LPUSH", queue, responseBody)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", queue, config.Conf.CommandTimeout)
	return
}

func respondSuccess(requestID string, result interface{}) error {
	return respond(requestID, &Response{
		Status: "success",
		Result: result,
	})
}

func respondError(requestID string, result interface{}) error {
	return respond(requestID, &Response{
		Status: "error",
		Result: result,
	})
}

func Listen(queue string) {
	conn := redisPool.Get()
	defer conn.Close()

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
			respondError(id, "Received a command without 'name'")
			continue
		}

		cmd := commands[name]
		if cmd == nil {
			respondError(id, "Unknown command")
			continue
		}

		log.Println("Processing command:", cmd.Name)
		result, err := cmd.handler(request.Params)
		if err != nil {
			respondError(id, err.Error())
		} else {
			respondSuccess(id, result)
		}
	}
}
