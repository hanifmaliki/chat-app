package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter server address: ")
	address, _ := reader.ReadString('\n')
	address = address[:len(address)-1]

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		err = conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
