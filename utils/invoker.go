package utils

import (
	"errors"
	"fmt"
	"log"
	"minerva/types"
	"os"
	"os/exec"
)

var commands = make(map[string]int, 0)

func InvokeCommand(action types.Action) error {
	if commands[action.Id] != 0 {
		log.Println("ERROR: action " + action.Id + " is already running (PID: " + fmt.Sprint(commands[action.Id]) + ")")
		return errors.New("ERROR: action " + action.Id + " is already running (PID: " + fmt.Sprint(commands[action.Id]) + ")")
	}

	cmd := exec.Command(action.Command[0], action.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Start()

	commands[action.Id] = cmd.Process.Pid

	return nil
}
