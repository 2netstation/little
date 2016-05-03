package main

import (
	"fmt"
	"log"
	"net/http"
	"websocketPrj/master"

	"code.google.com/p/go.net/websocket"
)

func main() {

	fmt.Println("Server start! ")
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/socket", websocket.Handler(master.MasterHandler))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
	fmt.Println("Server stop!")
}
