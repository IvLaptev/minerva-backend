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
		for _, command := range COMMANDS {
			if msg.Command == command {
				return msg, nil
			}
		}
		err = errors.New("wrong command name")
	}
	return msg, err
}
