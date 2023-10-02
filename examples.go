package main

import "github.com/gofiber/fiber/v2"

var examplesBases = []ExampleBase{
	{
		title:        "mouseover",
		templateName: "examples/color",
		binding:      fiber.Map{"Trigger": "mouseenter"},
		description:  "The box will fetch a new color from the server when you hover it",
	},
	{
		title:        "every 1s",
		templateName: "examples/color",
		binding:      fiber.Map{"Trigger": "every 1s"},
		description:  "The box will fetch a new color from the server every second",
	},
	{
		title:        "every 1s with fade",
		templateName: "examples/color",
		binding: fiber.Map{
			"Trigger": "every 1s",
			"Animate": true,
		},
		description: "The box will fetch a new color from the server and fade it in using CSS transitions",
	},
	{
		title:        "get on load",
		templateName: "examples/get",
		binding:      fiber.Map{"Trigger": "load"},
		description:  "Fetches a new message from the server when the page loads",
	},
	{
		title:        "get after delay",
		binding:      fiber.Map{"Trigger": "load delay:2s"},
		templateName: "examples/get",
		description:  "Fetches a new message from the server when the page loads after a 2 second delay",
	},
	{
		title:        "click to edit",
		templateName: "examples/contacts/initial",
		binding:      fiber.Map{},
		description:  "Sends form to the backend directly when click the Submit button and returns the server state",
	},
	{
		title:        "hx-indicator",
		templateName: "examples/indicator",
		binding:      fiber.Map{},
		description:  "Uses the 'hx-indicator' attribute to show a loading indicator and the 'hx-disabled-elt' attribute to disable the button while the request is in flight",
	},
	{
		title:        "click to load",
		templateName: "examples/click_to_load/table",
		binding: fiber.Map{
			"Rows": []int{1, 2},
		},
		description: "Click the button to load more rows from the server",
	},
	{
		title:        "open modal",
		templateName: "examples/modal",
		binding:      fiber.Map{},
		description:  "Will open a modal when you click the button",
	},
}
