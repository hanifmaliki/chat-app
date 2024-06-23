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
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter server address: (localhost:8080) ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)
	if address == "" {
		address = "localhost:8080"
	}

	// Health check
	resp, err := http.Get("http://" + address + "/healthz")
	if err != nil {
		log.Fatalf("Failed to check health: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read health check response: %v\n", err)
	}

	if string(body) != "ok" {
		log.Fatalf("Health check failed: %s\n", body)
	}

	// Prompt for login or register
	fmt.Print("Do you want to login or register? (login/register): ")
	action := ""
	for {
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		if action == "login" || action == "register" {
			break
		} else {
			log.Println("No action")
		}
	}

	// Read username and password
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	credentials := Credentials{
		Username: username,
		Password: password,
	}

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v\n", err)
	}

	// Send POST request to login or register endpoint
	var endpoint string
	if action == "login" {
		endpoint = "http://" + address + "/login"
	} else if action == "register" {
		endpoint = "http://" + address + "/register"
	} else {
		log.Fatalf("Invalid action: %s\n", action)
	}

	resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to %s: %v\n", action, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("%s failed: %s\n", action, body)
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/user/"+username+"/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	for {
		fmt.Print("Options: [1] Create Chat Room, [2] Select Chat Room, [3] Direct Message, [4] Send Message\nEnter option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			createChatRoom(conn, chatRoomName)

		case "2":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			selectChatRoom(conn, chatRoomName)

		case "3":
			fmt.Print("Enter username to direct message: ")
			targetUsername, _ := reader.ReadString('\n')
			targetUsername = strings.TrimSpace(targetUsername)

			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			sendDirectMessage(conn, targetUsername, message)

		case "4":
			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			err = conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}

		default:
			fmt.Println("Invalid option")
		}
	}
}

func createChatRoom(conn *websocket.Conn, chatRoomName string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("create:"+chatRoomName))
	if err != nil {
		log.Println("Failed to create chat room:", err)
	}
}

func selectChatRoom(conn *websocket.Conn, chatRoomName string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("select:"+chatRoomName))
	if err != nil {
		log.Println("Failed to select chat room:", err)
	}
}

func sendDirectMessage(conn *websocket.Conn, username string, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("dm:"+username+":"+message))
	if err != nil {
		log.Println("Failed to send direct message:", err)
	}
}
