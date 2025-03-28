package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Frontend origin
		AllowMethods: "GET,POST,PUT,DELETE",   // Allowed HTTP methods
		AllowHeaders: "Content-Type",          // Allowed headers
	}))

	// Routes
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Ping Working")
	})

	app.Get("/ws/:roomID", websocket.New(handleWebSocket))
	//What it does: handles connection to room socket/ exchange of messages
	//Expects: roomid
	//Returns:-

	app.Get("/startGame/:roomID/:maxrounds", handleStartGame)
	//What it does: calls startgame function, broadcasts "Game started" and initializes the story map
	//Expects: roomID, number of maxrounds
	//Returns: -

	app.Post("/submitline/:roomID", handleSubmitStory)
	//What it does: gets the story, adds it to the map, checks for
	//round-end condition, and calls next round if round ends
	//Frontend needs to call getall story if all the rounds are done,
	//otherwise the next prompts will be broadcast

	//Expects: roomID, jSON in the body with fields: "username", "message"
	//Returns: Confirmation of the message receival

	//REDUNDANT as servenextround broadcasts the promptmap anyways
	app.Get("/getNextPrompt/:roomID", getNextPrompt)
	//What it does: -
	//Expects: -
	//Returns: -

	app.Get("/check-username/:roomID/:username", handleRoomCheck)
	//What it does: Checks if username already exists in that room
	//Expects: roomid, username
	//Returns: error or name available as a JSON return

	app.Get("/get-final-stories/:roomID", getStory)
	//What it does: gets final stories map
	//Expects: roomID
	//Returns: final stories map

	// Start server
	log.Fatal(app.Listen(":3000"))
}
