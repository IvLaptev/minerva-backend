package main

import (
	"fmt"
	"minerva/handlers"
	"minerva/utils"
	"net/http"
)

func main() {
	err := utils.ReadConfig()
	if err != nil {
		fmt.Println("Can't read configuration file")
		fmt.Println(err)
	}

	config := utils.GetConfig()

	handlers.SetDefaultActions(config.Actions)

	if config.Service.Master {
		go utils.SetupSlaves(config.Service.Slaves, config.Service.Host)
	}

	// Обработка подключений, для управления системой
	http.HandleFunc("/control", handlers.WsMasterController)
	http.HandleFunc("/master", handlers.HandleMaster)

	// Запуск сервера
	http.ListenAndServe(":"+config.Service.Port, nil)
}
