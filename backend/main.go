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

	app.Get("/startGame/:roomID/:maxrounds", handleStartGame)

	app.Post("/submitline/:roomID", handleSubmitStory)

	app.Get("/getNextPrompt/:roomID", getNextPrompt)

	app.Get("/check-username/:roomID/:username", handleRoomCheck)

	app.Get("/get-final-stories", getStory)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
