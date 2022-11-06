package types

import (
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type Action struct {
	Title       string
	Description string
	Id          string
	Connection  *websocket.Conn
}

func GetActions(data []interface{}) ([]Action, error) {
	var actions []Action = make([]Action, 0)
	for i := 0; i < len(data); i++ {
		var action Action
		err := mapstructure.Decode(data[i], &action)
		if err != nil {
			return actions, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}
