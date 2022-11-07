package handlers

import (
	"encoding/json"
	"log"
	"minerva/types"
	"minerva/utils"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func HandleMaster(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	log.Println("MASTER FOUND:", addr)

	// Подключение к мастеру
	u := url.URL{Scheme: "ws", Host: addr, Path: "/control"}
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer connection.Close()

	// Отправка возможных для выполнения задач
	actions := make([]interface{}, 0)
	mapstructure.Decode(utils.GetConfig().Actions, &actions)
	msg, _ := json.Marshal(types.Message{Command: "actions.set", Body: actions})
	connection.WriteMessage(websocket.TextMessage, msg)

	// Сохранение соединения между сервисами
	WsSlaveHandler(connection)
}
