package main

import (
	"minerva/handlers"
	"net/http"
)

func main() {
	// Обработка подключения, для управления системой
	http.HandleFunc("/control", handlers.WsController)

	// Запуск сервера
	http.ListenAndServe(":8080", nil)
}
