package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"minerva/types"
	"minerva/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
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

	log.Println(actions)
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
		case types.COMMANDS[0]: // actions.set
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

		case types.COMMANDS[1]: // actions.get
			var resp_actions = make([]types.ResponceAction, 0)
			for _, action := range actions {
				resp_actions = append(resp_actions, action.ToResponseModel())

			}

			_SendMessage(connection, "actions.set", resp_actions)

		case types.COMMANDS[2]: // action.start
			action_id := msg.Body[0].(string)

			action := actions[action_id]
			if action.Connection != nil {
				action.Connection.WriteMessage(websocket.TextMessage, message)
			} else {
				log.Println("INVOKE ACTION:", action_id, action.Title)
				stdout := utils.InvokeCommand(action)

				// Ожидание окончания выполнения команды и отправка сообщения пользователю об этом
				go func() {
					// cmd.Wait()

					// Построчная отправка логов
					in := bufio.NewScanner(stdout)

					for in.Scan() {
						fmt.Println(in.Text())
						fmt.Println(types.ActionLog{Id: action_id, Line: in.Text()})
						_SendMessage(connection, "action.logs", []types.ActionLog{
							{Id: action_id, Line: in.Text()},
						})
					}

					// Обработка остановки выполнения команды
					utils.StopCommand(action)
					_SendMessage(connection, "action.stop", []string{action_id})
				}()
			}

		case types.COMMANDS[3]: // action.stop
			action_id := msg.Body[0].(string)

			action := actions[action_id]
			if action.Connection != nil {
				action.Connection.WriteMessage(websocket.TextMessage, message)
			} else {
				log.Println("STOP ACTION:", action_id, action.Title)
				utils.StopCommand(action)
			}
		}
	}
}

func _SendMessage(conn *websocket.Conn, command string, body interface{}) {
	resp_msg := types.Message{Command: command, Body: nil}
	mapstructure.Decode(body, &resp_msg.Body)
	resp, _ := json.Marshal(resp_msg)
	conn.WriteMessage(websocket.TextMessage, resp)
}
