package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"slices"
	"strconv"
	"time"
)

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

func colorHandler(c *fiber.Ctx) error {
	current := c.Query("current", "")
	trigger := c.Query("trigger", "")
	animate := c.Query("animate", "false")
	foo := slices.Index(colors, current)

	selectedIndex := (foo + 1) % len(colors)
	return c.Render("examples/color", fiber.Map{
		"Color":   colors[selectedIndex],
		"Trigger": trigger,
		"Animate": animate == "true",
	})
}

func getHandler(c *fiber.Ctx) error {
	return c.SendString("Hello from server")
}

var serverVersion = strconv.FormatInt(time.Now().Unix(), 10)

func reloadHandler(c *fiber.Ctx) error {
	clientVersion := c.Query("timestamp", "")
	if clientVersion != serverVersion {
		c.Set("HX-Refresh", "true")
	}
	return c.SendString("")
}

func sseHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		var i int
		msg := fmt.Sprintf("%d - the 2time is %v", i, time.Now())
		fmt.Fprintf(w, "event: TriggerReload\n")
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
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
}
