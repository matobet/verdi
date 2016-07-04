package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/env"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Ws struct {
	env.Commander
}

func NewHandler(cmd env.Commander) http.Handler {
	ws := &Ws{cmd}
	return sockjs.NewHandler("/ws", sockjs.DefaultOptions, ws.Handler)
}

func (ws *Ws) Handler(session sockjs.Session) {
	for {
		if msg, err := session.Recv(); err == nil {
			log.Printf("ws: received: %v", msg)
			message := map[string]interface{}{}
			err = json.Unmarshal([]byte(msg), &message)
			if err != nil {
				respondError(session, err.Error())
				continue
			}
			command, ok := message["cmd"]
			if !ok {
				respondError(session, "`cmd` field must be specified!")
				continue
			}
			commandName, ok := command.(string)
			if !ok {
				respondError(session, "`cmd` must be a string")
				continue
			}
			params, ok := message["params"]
			if !ok {
				respondError(session, "`params` field must be specified!")
				continue
			}
			paramsMap, ok := params.(map[string]interface{})
			if !ok {
				respondError(session, "`params` must be an object!")
				continue
			}
			result, err := ws.Run(commandName, paramsMap)
			if err != nil {
				respondError(session, err.Error())
				continue
			}
			response, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			session.Send(string(response))
		} else {
			log.Println(err)
			break
		}
	}
}

func respondError(session sockjs.Session, message string) {
	result := cmd.Response{
		Status: "error",
		Result: message,
	}
	response, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	session.Send(string(response))
}
