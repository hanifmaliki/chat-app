package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hanifmaliki/chat-app/internal/model"
)

func main() {
	address := promptInput("Enter server address: (localhost:8080) ", "localhost:8080")

	if !healthCheck(address) {
		log.Fatalf("Health check failed")
	}

	for {
		action := promptInput("Do you want to login or register? (login/register): ", "")
		if action != "login" && action != "register" {
			log.Println("Invalid action, please enter 'login' or 'register'.")
			continue
		}

		username := promptInput("Enter username: ", "")
		password := promptInput("Enter password: ", "")

		if !performAuth(address, action, username, password) {
			if action == "register" {
				continue
			} else {
				log.Fatalf("%s failed", action)
			}
		}

		conn := connectWebSocket(address, username)
		defer conn.Close()

		go readMessages(conn)
		mainLoop(conn)
	}
}

func promptInput(prompt, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func healthCheck(address string) bool {
	resp, err := http.Get("http://" + address + "/health")
	if err != nil {
		log.Printf("Failed to check health: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read health check response: %v\n", err)
		return false
	}

	return string(body) == "ok"
}

func performAuth(address, action, username, password string) bool {
	credentials := model.Credential{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v\n", err)
	}

	endpoint := "http://" + address + "/" + action
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to %s: %v\n", action, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("%s failed: %s\n", action, body)
		return false
	}

	return true
}

func connectWebSocket(address, username string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/chat/ws?username="+username, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v\n", err)
	}
	return conn
}

func readMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v\n", err)
			return
		}
		fmt.Println(string(message))
	}
}

func mainLoop(conn *websocket.Conn) {
	for {
		fmt.Print("Options:\n[1] Direct Message\n[2] Create Chat Room\n[3] Join Chat Room\n[4] Leave Chat Room\n[5] Send Message to Chat Room\nEnter option: ")
		option := promptInput("", "")

		switch option {
		case "1":
			targetUsername := promptInput("Enter username to direct message: ", "")
			message := promptInput("Enter message: ", "")
			sendMessage(conn, "dm:"+targetUsername+":"+message)

		case "2":
			chatRoomName := promptInput("Enter chat room name: ", "")
			sendMessage(conn, "room:create:"+chatRoomName)

		case "3":
			chatRoomName := promptInput("Enter chat room name: ", "")
			sendMessage(conn, "room:join:"+chatRoomName)

		case "4":
			chatRoomName := promptInput("Enter chat room name: ", "")
			sendMessage(conn, "room:leave:"+chatRoomName)

		case "5":
			chatRoomName := promptInput("Enter chat room name: ", "")
			message := promptInput("Enter message: ", "")
			sendMessage(conn, "room:broadcast:"+chatRoomName+":"+message)

		default:
			fmt.Println("Invalid option")
		}
	}
}

func sendMessage(conn *websocket.Conn, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
	}
}
