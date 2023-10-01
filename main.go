package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/codecat/melody"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"
	gohtml "html"
	"html/template"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
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
		templateName: "examples/get-color",
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
		templateName: "examples/get-color",
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
		templateName: "examples/get-load",
		binding:      fiber.Map{},
		description:  "This is a simple example",
		attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load"`,
		},
	},
	{
		title:        "Example 2",
		binding:      fiber.Map{},
		templateName: "examples/get-load-delay",
		description:  "This is a simple example",
		attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load delay:2s"`,
		},
	},
}

var colors = []string{
	"bg-gray-100",
	"bg-red-200",
	"bg-yellow-300",
	"bg-green-400",
	"bg-blue-500",
	"bg-indigo-600",
	"bg-purple-700",
	"bg-pink-800",
}

var isDev = os.Getenv("ENV") == "dev"

func main() {
	serverVersion := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(serverVersion)

	// Create a new engine
	engine := html.New("static/views", ".html")
	engine.Debug(isDev)
	engine.Reload(isDev)

	engine.AddFunc("escape", func(s string) string {
		return gohtml.EscapeString(s)
	})

	examples := generateHtmxExamples(engine)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/styles", "./static/styles")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":         "Hello, HTMX!",
			"ServerVersion": serverVersion,
			"Examples":      examples,
		}, "layouts/main")
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		return c.SendString("Hello from server")
	})
	app.Get("/reload", func(c *fiber.Ctx) error {
		clientVersion := c.Query("timestamp", "")
		if clientVersion != serverVersion {
			c.Set("HX-Refresh", "true")
		}
		return c.SendString("")
	})
	app.Get("/color", func(c *fiber.Ctx) error {
		current := c.Query("current", "")
		trigger := c.Query("trigger", "")
		foo := slices.Index(colors, current)

		selectedIndex := (foo + 1) % len(colors)
		return c.Render("examples/get-color", fiber.Map{
			"Color":   colors[selectedIndex],
			"Trigger": trigger,
		})
	})

	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			var i int
			msg := fmt.Sprintf("%d - the 2time is %v", i, time.Now())
			fmt.Fprintf(w, "event: TriggerReload\n")
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			fmt.Println(msg + "\n")
			err := w.Flush()
			for {
				i++
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
					break
				}
				// TODO: lock here forever, instead?
				time.Sleep(1 * time.Second)
			}
		}))

		return nil
	})

	// Websocket stuff
	m := melody.New()

	app.Get("/ws", func(c *fiber.Ctx) error {
		return m.HandleRequest(c.Context())
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var in WsMessage
		err := json.Unmarshal(msg, &in)
		if err != nil {
			fmt.Println("Error parsing message", err)
			return
		}

		//msg2, err := json.Marshal(in)
		//if err != nil {
		//	fmt.Println("Error marshalling message", err)
		//	return
		//}

		m.Broadcast([]byte(fmt.Sprintf("<div hx-swap-oob='beforeend:#messages'><p><b>{username}</b>: %s</p></div><div hx-swap-oob='beforeend:#messages2'><p>%s,%s</p></div>", in.ChatMessage, in.ChatMessage, in.ChatMessage)))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(app.Listen(":" + port))
}

type WsMessage struct {
	ChatMessage string `json:"chat_message"`
}

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
