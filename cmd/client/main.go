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

	"github.com/hanifmaliki/chat-app/internal/model"

	"github.com/gorilla/websocket"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter server address: (localhost:8080) ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)
	if address == "" {
		address = "localhost:8080"
	}

	// Health check
	// resp, err := http.Get("http://" + address + "/healthz")
	// if err != nil {
	// 	log.Fatalf("Failed to check health: %v\n", err)
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Failed to read health check response: %v\n", err)
	// }

	// if string(body) != "ok" {
	// 	log.Fatalf("Health check failed: %s\n", body)
	// }

LOGIN:

	// Prompt for login or register
	fmt.Print("Do you want to login or register? (login/register): ")
	action := ""
	for {
		action, _ = reader.ReadString('\n')
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

	credentials := model.Credential{
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

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to %s: %v\n", action, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("%s failed: %s\n", action, body)
	}

	if action == "register" {
		goto LOGIN
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+address+"/chat/ws?username="+username, nil)
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
			fmt.Println(string(message))
		}
	}()

	for {
		fmt.Print("Options:\n[1] Direct Message\n[2] Create Chat Room\n[3] Join Chat Room\n[4] Leave Chat Room\n[5] Send Message to Chat Room\nEnter option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("Enter username to direct message: ")
			targetUsername, _ := reader.ReadString('\n')
			targetUsername = strings.TrimSpace(targetUsername)

			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			sendDirectMessage(conn, targetUsername, message)

		case "2":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			createChatRoom(conn, chatRoomName)

		case "3":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			joinChatRoom(conn, chatRoomName)

		case "4":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			leaveChatRoom(conn, chatRoomName)

		case "5":
			fmt.Print("Enter chat room name: ")
			chatRoomName, _ := reader.ReadString('\n')
			chatRoomName = strings.TrimSpace(chatRoomName)

			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			messageChatRoom(conn, chatRoomName, message)

		default:
			fmt.Println("Invalid option")
		}
	}
}

func sendDirectMessage(conn *websocket.Conn, username string, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("dm:"+username+":"+message))
	if err != nil {
		log.Println("Failed to send direct message:", err)
	}
}

func createChatRoom(conn *websocket.Conn, chatRoomName string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("room:create:"+chatRoomName))
	if err != nil {
		log.Println("Failed to create chat room:", err)
	}
}

func joinChatRoom(conn *websocket.Conn, chatRoomName string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("room:join:"+chatRoomName))
	if err != nil {
		log.Println("Failed to join chat room:", err)
	}
}

func leaveChatRoom(conn *websocket.Conn, chatRoomName string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("room:leave:"+chatRoomName))
	if err != nil {
		log.Println("Failed to leave chat room:", err)
	}
}

func messageChatRoom(conn *websocket.Conn, chatRoomName string, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("room:broadcast:"+chatRoomName+":"+message))
	if err != nil {
		log.Println("Failed to send message to chat room:", err)
	}
}
