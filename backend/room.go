package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/s-bhatla/chatroom/utils"

	"github.com/gofiber/websocket/v2"
)

type StoryLine struct {
	user  string
	story string
}

type Room struct {
	mu          sync.Mutex
	connections map[*websocket.Conn]string
	maxrounds   uint8
	round       uint8
	Stories     map[string][]StoryLine //username storystarter mapped to array of complete story (which will be appended every round)
	//Userlist must be changed everytime user disconnects
	userlist []string //usernames list
}

var rooms = make(map[string]*Room)
var roomsMu sync.Mutex

func getOrCreateRoom(roomID string) *Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if room, exists := rooms[roomID]; exists {
		return room
	}

	room := &Room{
		connections: make(map[*websocket.Conn]string),
	}
	rooms[roomID] = room

	return room
}

func (r *Room) startGame(maxrounds uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//When all the players are in the room and the game is to be started.

	//set rounds of the room
	r.maxrounds = maxrounds

	//Change this to add the null or sth
	for _, user := range r.userlist {
		newStoryline := StoryLine{user: user, story: ""}
		r.Stories[user] = make([]StoryLine, 0)
		r.Stories[user] = append(r.Stories[user], newStoryline)
	}
	//Just send it to frontend and get the stories responses
	r.round = 1
	r.broadcast("Game started", nil)
}

func (r *Room) addConnection(username string, conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, alreadyname := range r.connections {
		if alreadyname == username {
			return errors.New("username already exists")
		}
	}
	r.connections[conn] = username
	return nil
}

func (r *Room) removeConnection(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.connections, conn)

	// Check if the game should end due to disconnection
	if len(r.connections) < len(r.Stories) {
		r.broadcast("Error: Game ended due to player disconnection", nil)
		r.endGame()
	}
}

func (r *Room) serveNextRound() {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows := len(r.Stories)

	matrix := make([][]string, 0, rows)

	for _, storyline := range r.Stories {
		users := make([]string, 0, len(storyline))

		for _, line := range storyline {
			users = append(users, line.user)
		}

		matrix = append(matrix, users)
	}

	//take a slice of the first list to send to fxn as userlist
	var matrixSlice [][]string
	if len(matrix[0]) > len(matrix) {
		// Take modulus % and use that slice to send to fxn
		matrixSlice = matrix[:len(matrix)%len(matrix[0])]
	} else {
		// Slice to send = matrix
		matrixSlice = matrix
	}

	// Call get new col and get slice of new column
	newslice := utils.GetNewCol(matrixSlice, []string{}, r.userlist)

	//Use the created and sent matrrix as a reference
	// for which user goes where as per the output

	promptMap := make(map[string]string) //A map which defines which user will get which which prompt
	// basically just the previous storyline

	for i, user := range matrixSlice[0] {
		promptMap[newslice[i]] = r.Stories[user][len(r.Stories[user])-1].story

		newStoryline := StoryLine{user: newslice[i], story: ""}
		r.Stories[user] = append(r.Stories[user], newStoryline)
	}

	fmt.Println(promptMap)
	// Broadcast the updated stories to all users in the room
	r.broadcast("New round started", nil)

	jsonPromptMap, err := json.Marshal(promptMap)
	if err != nil {
		log.Println(err)
		return
	}

	// Broadcast JSON promptMap to all connected users
	r.broadcast(string(jsonPromptMap), nil)

}

func (r *Room) broadcast(message string, sender *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for conn := range r.connections {
		// Send message to everyone but the sender
		if conn != sender {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error broadcasting to client %v", err)
			}
		}
	}
}

// and that realtime updation shi (submitted/not submitted tick feature) on submitted (POST req) add onto handleSubmitStory......
// Can simply use the broadcast function for that no worries - moment all submitted are true (make a check everytime one goes true)
//call the function to serve new lines - (gotta do in the frontend...)

//Combine all to a flow and get frontend
