package main

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"math/big"
)
//even tho re are crypto random not sudo random we still check and keep record of what we used already 
var generatedNumbers = make(map[string]bool)

func generateUniqueRandomNumber() *big.Int {
	for {
		// Generate a large random number (adjust the bit size as needed)
		n, err := rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil))
		if err != nil {
			panic(err)
		}

		if exists := generatedNumbers[n.String()]; !exists {
			// If not, mark it as generated and return it
			generatedNumbers[n.String()] = true
			return n
		} else {
			// If it exists, loop until a unique number is generated
			for exists {
				n, err = rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil))
				if err != nil {
					panic(err) // Handle the error appropriately
				}
				if _, exists = generatedNumbers[n.String()]; !exists {
					generatedNumbers[n.String()] = true
					return n
				}
			}
		}

		// If it exists, continue the loop to generate a new number
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accept all connections as document said no security or restriction what so ever fun
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	answer:= generateUniqueRandomNumber()

	// Simple loop to keep the connection open
	for {
		// Read message from browser/client etc...
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Log the received message
		log.Printf("Received: %s\n", msg)

		bytesofbigint := answer.Bytes() 
		log.Printf("we responded with: %+v\n", answer)

		// response to client I could send back the string verion via websocket.textmessage but 
		// it was requested that it should always be *big.int even tho your typical client would show this as something 
		//nonsense like xQ�c�ҽ2g�)U� 
		//anyway the human readble option is log.Printf("we responded with: %+v\n", answer)
		if err := ws.WriteMessage(websocket.BinaryMessage,  bytesofbigint ); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	// Configure WebSocket route lets not make it complicated
	http.HandleFunc("/ws", handleConnections)

	// Start the server on localhost port 8080 because this randomiser server is so cool
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Listen And  Serve: ", err)
	}
}
