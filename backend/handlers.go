package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (r *Room) isUsernameAvailable(username string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existingUsername := range r.connections {
		if existingUsername == username {
			return false
		}
	}
	return true
}

func getStory(c *fiber.Ctx) error {
	roomID := c.Params("roomID")

	if room, exists := rooms[roomID]; exists {
		return c.Status(fiber.StatusOK).JSON(room.Stories)
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "room does not exist.",
		})
	}
}

func handleRoomCheck(c *fiber.Ctx) error {
	roomID := c.Params("roomID")
	username := c.Params("username")

	if room, exists := rooms[roomID]; exists {
		if room.isUsernameAvailable(username) {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"available": true,
				"message":   "Username is available",
			})
		} else {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"available": false,
				"message":   "Username is already taken",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"available": true,
		"message":   "Room not found, creating new room.",
	})
}

func handleWebSocket(c *websocket.Conn) {
	roomID := c.Params("roomID")

	// Ensure the room exists, if not - create
	room := getOrCreateRoom(roomID)

	//Gotta change this
	// username := "RandomUser"
	var username string

	if err := c.ReadJSON(&username); err != nil || username == "" {
		log.Printf("Invalid username %v from client in room %s: %v", username, roomID, err)
		c.Close()
		return
	}

	if err := room.addConnection(username, c); err != nil {
		log.Printf("Error adding user %s to the room %s: %v", username, roomID, err)
		c.WriteJSON(map[string]string{"error": err.Error()})
		c.Close()
		return
	}
	log.Printf("Client connected to room: %s", roomID)

	// Handling messages from client
	defer func() {
		room.removeConnection(c)
		log.Printf("Client disconnected from room %s", roomID)
		c.Close()
	}()

	for {
		var msg string
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("Error while reading the message from room %s: %v", roomID, err)
			break
		}
		log.Printf("Message from room %s: %s", roomID, msg)

		room.broadcast(msg, c)
	}
}

type storyContribution struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func getNextPrompt(ctx *fiber.Ctx) error {
	//get the username
	//Call the fxn in the room, (to be written)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "And the story continued...",
	})
}

// change function
func handleSubmitStory(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	room := getOrCreateRoom(roomID)
	// Parse the request body into a MessageRequest struct
	request := new(storyContribution)
	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	print(request)
	//iterate over stories to find a Stoyline type with the same username as our boy and empty story (we'll set it that way before serving)
	for _, value := range room.Stories {
		for _, storyline := range value {
			if storyline.user == request.Username && storyline.story == "" {
				storyline.story = request.Message
			}
		}
	}

	// Respond with the parsed data
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message received successfully",
		"data":    request,
	})
}
