package main

import (
	"fmt"
	"minerva/handlers"
	"minerva/types"
	"net/http"
)

func main() {
	err := types.ReadConfig()
	if err != nil {
		fmt.Println("Can't read configuration file")
		fmt.Println(err)
	}

	config := types.GetConfig()

	// Обработка подключения, для управления системой
	http.HandleFunc("/control", handlers.WsController)

	// Запуск сервера
	http.ListenAndServe(":"+config.Service.Port, nil)
}
