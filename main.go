package main

import (
	"bytes"
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	gohtml "html"
	"html/template"
	"log"
	"os"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(app.Listen(":" + port))
}
