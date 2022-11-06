package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Message struct {
	Command string
	Body    []interface{}
}

func GetMessageFromBytes(message []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(message, &msg)

	if err != nil {
		fmt.Println(err)
	} else {
		if msg.Command == COMMANDS[0] || msg.Command == COMMANDS[1] || msg.Command != COMMANDS[2] {
			return msg, nil
		} else {
			err = errors.New("wrong command name")
		}
	}
	return msg, err
}
