package utils

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func SetupSlaves(addrs []string, master string) {
	for _, addr := range addrs {
		go find_save(addr, master)
	}
}

func find_save(addr string, master string) {
	for {
		time.Sleep(5 * time.Second)
		r, err := http.Get(addr + "/master?address=" + master)

		if err == nil && strings.Split(r.Status, " ")[0] == "200" {
			break
		} else {
			log.Println("ERROR: Can't connect to", addr)
		}
	}
}
