package main

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	templts "github.com/magnuswahlstrand/htmx-experiments/components"
	"log"
	"os"
)

var isDev = os.Getenv("ENV") == "dev"

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		w := templts.Page(serverVersion)
		c.Set("Content-Type", "text/html")
		return w.Render(c.Context(), c.Response().BodyWriter())
	})
	app.Static("/", "./static")
	app.Get("/get", getHandler)
	app.Get("/reload", reloadHandler)
	app.Get("/color", colorHandler)
	app.Get("/sse", sseHandler)
	app.Get("/track", trackHandler)
	app.Post("/slow", slowHandler)
	contacts := app.Group("/contacts")
	contacts.Put("/1", contactsUpdatePutHandler)
	contacts.Get("/1", contactGetHandler)
	contacts.Get("/1/edit", contactEditGetHandler)
	app.Get("/click_to_load", clickToLoadHandler)
	app.Get("/modal", modalHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(app.Listen(":" + port))
}
