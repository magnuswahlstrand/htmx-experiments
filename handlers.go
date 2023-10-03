package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	templts "github.com/magnuswahlstrand/htmx-experiments/components"
	"github.com/magnuswahlstrand/htmx-experiments/types"
	"github.com/valyala/fasthttp"
	"net/http"
	"slices"
	"strconv"
	"sync"
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
	currentIndex := slices.Index(colors, current)
	color := colors[(currentIndex+1)%len(colors)]

	w := templts.Color(trigger, color, animate == "true")
	return w.Render(c.Context(), c.Response().BodyWriter())
}

func trackHandler(c *fiber.Ctx) error {
	currentState, err := strconv.Atoi(c.Query("state", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	w := templts.ExampleTrack(currentState)
	return w.Render(c.Context(), c.Response().BodyWriter())
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
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				break
			}
			// TODO: lock here forever, instead?
			time.Sleep(1 * time.Second)
		}
	}))

	return nil
}

func slowHandler(ctx *fiber.Ctx) error {
	time.Sleep(1 * time.Second)
	return ctx.SendStatus(http.StatusNoContent)
}

func clickToLoadHandler(c *fiber.Ctx) error {
	time.Sleep(100 * time.Millisecond)

	pageStr := c.Query("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	agentID := 2*page + 1
	w := templts.ClickToLoadRows([]int{agentID, agentID + 1}, page+1)
	return w.Render(c.Context(), c.Response().BodyWriter())
}

func modalHandler(c *fiber.Ctx) error {
	w := templts.Modal()
	return w.Render(c.Context(), c.Response().BodyWriter())
}

var contactMu = &sync.Mutex{}
var contact = types.Contact{
	Name:  "Magnus",
	Email: "magnus@mail.com",
}

func contactGetHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()
	w := templts.ContactForm(contact, false)
	return w.Render(c.Context(), c.Response().BodyWriter())
}

func contactEditGetHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()

	w := templts.ContactForm(contact, true)
	return w.Render(c.Context(), c.Response().BodyWriter())
}

func contactsUpdatePutHandler(c *fiber.Ctx) error {
	contactMu.Lock()
	defer contactMu.Unlock()

	var update types.Contact
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	contact.Email = update.Email
	contact.Name = update.Name
	w := templts.ContactForm(contact, false)
	return w.Render(c.Context(), c.Response().BodyWriter())
}
