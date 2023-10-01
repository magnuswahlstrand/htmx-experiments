package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	gohtml "html"
	"html/template"
	"log"
	"os"
	"sync"
)

type ExampleBase struct {
	title        string
	templateName string
	binding      fiber.Map
	description  string
	attributes   []string
}

type Example struct {
	Title       string
	Component   template.HTML
	Description string
	Attributes  []string
}

var examplesBases = []ExampleBase{
	{
		title:        "click to edit",
		templateName: "examples/contacts/get",
		binding:      contact.Bindings(false),
		description:  "This is a simple example",
		attributes:   []string{},
	},
	{
		title:        "trigger: mouseover",
		templateName: "examples/color",
		binding: fiber.Map{
			"Trigger": "mouseenter",
		},
		description: "This is a simple example",
		attributes: []string{
			`hx-get="/color"`,
			`hx-trigger="mouseenter"`,
		},
	},
	{
		title:        "trigger: every 1s",
		templateName: "examples/color",
		binding: fiber.Map{
			"Trigger": "every 1s",
		},
		description: "This is a simple example",
		attributes: []string{
			`hx-get="/color"`,
			`hx-trigger="every 1s"`,
		},
	},
	{
		title:        "Example 1",
		templateName: "examples/get",
		binding: fiber.Map{
			"Trigger": "load",
		},
		description: "This is a simple example",
		attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load"`,
		},
	},
	{
		title: "Example 2",
		binding: fiber.Map{
			"Trigger": "load delay:2s",
		},
		templateName: "examples/get",
		description:  "This is a simple example",
		attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load delay:2s"`,
		},
	},
}

var isDev = os.Getenv("ENV") == "dev"

func generateHtmxExamples(engine *html.Engine) []Example {
	var examples []Example
	for _, v := range examplesBases {
		v := v
		var buffer bytes.Buffer
		err := engine.Render(&buffer, v.templateName, v.binding)
		if err != nil {
			log.Fatal("Failed to render template", err)
		}

		examples = append(examples, Example{
			Title:       v.title,
			Component:   template.HTML(buffer.String()),
			Description: v.description,
			Attributes:  v.attributes,
		})
	}
	return examples
}
func main() {

	// Create a new engine
	engine := html.New("static/views", ".html")
	engine.Debug(isDev)
	engine.Reload(isDev)

	engine.AddFunc("escape", func(s string) string {
		return gohtml.EscapeString(s)
	})

	examples := generateHtmxExamples(engine)
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":         "Hello, HTMX!",
			"ServerVersion": serverVersion,
			"Examples":      examples,
		}, "layouts/main")
	})
	app.Static("/styles", "./static/styles")
	app.Get("/get", getHandler)
	app.Get("/reload", reloadHandler)
	app.Get("/color", colorHandler)
	app.Get("/sse", sseHandler)
	contacts := app.Group("/contacts")
	contacts.Put("/1", contactsUpdatePutHandler)
	contacts.Get("/1", contactGetHandler)
	contacts.Get("/1/edit", contactEditGetHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(app.Listen(":" + port))
}

type Contact struct {
	Name  string
	Email string
}

func (c Contact) Bindings(edit bool) fiber.Map {
	return fiber.Map{
		"Name":  c.Name,
		"Email": c.Email,
		"Edit":  edit,
	}
}

var contactMu = &sync.Mutex{}
var contact = Contact{
	Name:  "Magnus",
	Email: "magnus@mail.com",
}

func contactGetHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()
	return c.Render("examples/contacts/get", contact.Bindings(false))
}

func contactEditGetHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()

	return c.Render("examples/contacts/edit", contact.Bindings(true))
}

func contactsUpdatePutHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()

	var update Contact
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	contact.Email = update.Email
	contact.Name = update.Name
	fmt.Println("contact", contact)
	return c.Render("examples/contacts/get", contact.Bindings(false))
}
