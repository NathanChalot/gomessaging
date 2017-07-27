package main

import (
  "log"
  "net/http"
  "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // connected clients
var broadcast = make(chan Message)           // broadcast channel
// Configure the upgrader
var upgrader = websocket.Upgrader{}
// Define our message object
type Message struct {
  Type     string `json:"type"`
  Email    string `json:"email"`
  Username string `json:"username"`
  Message  string `json:"message"`
  Destination string `json:"destination"`
}

func main() {
  // Create a simple file server
  fs := http.FileServer(http.Dir("public"))
  http.Handle("/", fs)
  // Configure websocket route
  http.HandleFunc("/ws", handleConnections)
  // Start listening for incoming chat messages
  go handleMessages()
  // Start the server on localhost port 8000 and log any errors
  log.Println("http server started on :8000")
  err := http.ListenAndServe(":8000", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
  // Upgrade initial GET request to a websocket
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  // Make sure we close the connection when the function returns
  defer ws.Close()
  // Register our new client
  clients[ws] = ""
  for {
    var msg Message
    // Read in a new message as JSON and map it to a Message object
    err := ws.ReadJSON(&msg)
    if err != nil {
      log.Printf("error: %v", err)
      delete(clients, ws)
      break
    }
    if msg.Type == "connect" {
      clients[ws] = msg.Username
      continue
    }
    // Send the newly received message to the broadcast channel
    broadcast <- msg
  }
}

func handleMessages() {
  for {
    // Grab the next message from the broadcast channel
    msg := <-broadcast
    // Send it out to every client that is currently connected
    for client, name := range clients {
      log.Printf("client = %v", name)
      log.Printf("destination = %v", msg.Destination)
      if (msg.Username == name || msg.Destination == name ||
        msg.Destination == "") {
        log.Printf("SENDING MESSAGE to %v", name)
        err := client.WriteJSON(msg)
        if err != nil {
          log.Printf("error: %v", err)
          client.Close()
          delete(clients, client)
        }
      }
    }
  }
}
