package utils

import (
	"errors"
	"fmt"
	"log"
	"minerva/types"
	"os"
	"os/exec"
	"syscall"
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
	log.Println("ACTION STARTED: ID: " + action.Id + "; PID: " + fmt.Sprint(commands[action.Id]))

	return nil
}

func StopCommand(action types.Action) error {
	pid := commands[action.Id]
	if pid == 0 {
		log.Println("ERROR: action " + action.Id + " is not running")
		return errors.New("ERROR: action " + action.Id + " is not running")
	}
	defer delete(commands, action.Id)

	syscall.Kill(pid, syscall.SIGTERM)
	log.Println("ACTION STOPPED: ID: " + action.Id + "; PID: " + fmt.Sprint(commands[action.Id]))

	return nil
}
