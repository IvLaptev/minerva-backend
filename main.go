package main

import (
	"fmt"
	"log"
	"minerva/types"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

// Соединения с клиентами
var clients map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
var actions map[string]types.Action = make(map[string]types.Action)

func main() {
	// Обработка подключения, для управления системой
	http.HandleFunc("/control", control)

	// Запуск сервера
	http.ListenAndServe(":8080", nil)
}

// Перевод подключения с http на ws
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func control(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()
	log.Println("CONNECTED:", connection.RemoteAddr())

	clients[connection] = true
	defer delete(clients, connection)

	// Обработка сообщений по ws
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			log.Println("DISCONNECTED:", connection.RemoteAddr())
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		var msg types.Message
		msg, err = types.GetMessageFromBytes(message)
		if err != nil {
			continue
		}

		switch msg.Command {
		case types.COMMANDS[0]: // set_actions
			new_actions, err := types.GetActions(msg.Body)
			if err != nil {
				continue
			}

			for i := 0; i < len(new_actions); i++ {
				actions[new_actions[i].Id] = new_actions[i]
				defer delete(actions, new_actions[i].Id)
			}
		}

		log.Println("FROM:", connection.RemoteAddr(), "MESSAGE:", string(message))
		fmt.Println(actions)
		connection.WriteMessage(websocket.TextMessage, message)

		go messageHandler(message)
	}

}

// Обработка сообщения, полученного по ws
func messageHandler(message []byte) {

	cmd := exec.Command("sh", "./long_task.sh")
	// cmd.Stdout = os.Stdout
	cmd.Start()
}

// switch msg.Command {
// 	case types.COMMANDS[0]: // set_actions
// 		var new_actions []types.Action = make([]types.Action, 0)
// 		for i := 0; i < len(msg.Body); i++ {
// 			var action types.Action
// 			mapstructure.Decode(msg.Body[i], &action)
// 			new_actions = append(new_actions, action)
// 		}

// 		actions[]

// 		// case COMMANDS[1]:
// 		// case COMMANDS[2]:
// 	}
