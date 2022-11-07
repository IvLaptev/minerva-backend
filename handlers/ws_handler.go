package handlers

import (
	"encoding/json"
	"log"
	"minerva/types"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

// Соединения с клиентами
var clients map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

// Задачи, выполнение которых можно запустить
var actions map[string]types.Action = make(map[string]types.Action)

// Перевод подключения с http на ws
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func SetDefaultActions(new_actions []types.Action) {

	for _, action := range new_actions {
		if actions[action.Id].Id != "" {
			log.Println("ERROR: Action with ID", action.Id, "already exists")
			continue
		}

		actions[action.Id] = action
	}
}

func WsMasterController(w http.ResponseWriter, r *http.Request) {
	// Установка соединения
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()
	log.Println("CONNECTED:", connection.RemoteAddr())

	clients[connection] = true
	defer delete(clients, connection)

	// Обработка сообщений соединения
	WsSlaveHandler(connection)
}

func WsSlaveHandler(connection *websocket.Conn) {
	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			log.Println("DISCONNECTED:", connection.RemoteAddr())
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}

		log.Println("FROM:", connection.RemoteAddr(), "MESSAGE:", string(message))

		// Обработка сообщения
		var msg types.Message
		msg, err = types.GetMessageFromBytes(message)
		if err != nil {
			continue
		}

		// Выполнение переданной команды
		switch msg.Command {
		case types.COMMANDS[0]: // set_actions
			new_actions, err := types.GetActions(msg.Body)
			if err != nil {
				continue
			}

			for i := 0; i < len(new_actions); i++ {
				if actions[new_actions[i].Id].Id != "" {
					log.Println("ERROR: Action with ID", new_actions[i].Id, "already exists")
					continue
				}

				new_actions[i].Connection = connection
				actions[new_actions[i].Id] = new_actions[i]
				defer delete(actions, new_actions[i].Id)
			}

		case types.COMMANDS[1]: // get_actions
			var resp_actions = make([]types.ResponceAction, 0)
			for _, action := range actions {
				resp_actions = append(resp_actions, types.ResponceAction{
					Id:          action.Id,
					Title:       action.Title,
					Description: action.Description,
				})

				resp, _ := json.Marshal(resp_actions)

				connection.WriteMessage(websocket.TextMessage, resp)
			}

		case types.COMMANDS[2]: // run_action
			action_id := msg.Body[0].(string)

			action := actions[action_id]
			if action.Connection != nil {
				action.Connection.WriteMessage(websocket.TextMessage, message)
			} else {
				log.Println("INVOKE ACTION")
			}
		}

		// Отправка ответа клиенту
		// connection.WriteMessage(websocket.TextMessage, message)

		go messageHandler(message)
	}
}

// Обработка сообщения, полученного по ws
func messageHandler(message []byte) {

	cmd := exec.Command("sh", "./long_task.sh")
	// cmd.Stdout = os.Stdout
	cmd.Start()
}
