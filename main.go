package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"
	gohtml "html"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

type Example struct {
	Title       string
	Component   template.HTML
	Description string
	Attributes  []string
}

//go:embed views/examples/get-load.html
var getLoad string

//go:embed views/examples/get-load-delay.html
var getLoadDelay string

var examples = []Example{
	{
		Title:       "Example 1",
		Component:   template.HTML(getLoad),
		Description: "This is a simple example",
		Attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load"`,
		},
	},
	{
		Title:       "Example 2",
		Component:   template.HTML(getLoadDelay),
		Description: "This is a simple example",
		Attributes: []string{
			`hx-get="/get"`,
			`hx-trigger="load delay:2s"`,
		},
	},
}

func main() {
	serverVersion := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(serverVersion)

	// Create a new engine
	engine := html.New("./views", ".html")
	engine.Debug(true)
	engine.Reload(true)

	engine.AddFunc("escape", func(s string) string {
		return gohtml.EscapeString(s)
	})

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := html.NewFileSystem(http.Dir("./views", ".html"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/dist", "./dist")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":         "Hello, HTMX!",
			"ServerVersion": serverVersion,
			"Examples":      examples,
		}, "layouts/main")
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		return c.SendString("Foo")
	})
	app.Get("/reload", func(c *fiber.Ctx) error {
		clientVersion := c.Query("timestamp", "")
		fmt.Println(clientVersion)
		if clientVersion != serverVersion {
			c.Set("HX-Refresh", "true")
		}
		return c.SendString("bb")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	log.Fatal(app.Listen(":" + port))
}
