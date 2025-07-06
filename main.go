package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var todos []Todo
var idCounter int

func main() {
	engine := html.New("./views", ".html")
	tenSecTimeout := 10 * time.Second
	oneMB := 1 << 20
	app := fiber.New(fiber.Config{
		Views: engine,
		BodyLimit: oneMB,
		Concurrency: 256 * 1024,
		IdleTimeout: tenSecTimeout,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Todos": todos,
		})
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		idCounter++
		todo := Todo{
			ID:   idCounter,
			Text: c.FormValue("text"),
			Done: false,
		}
		todos = append(todos, todo)
		return c.Redirect("/")
	})

	app.Post("/delete/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		for i, t := range todos {
			if t.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				break
			}
		}
		return c.Redirect("/")
	})

	app.Shutdown()
	app.ShutdownWithTimeout(tenSecTimeout)
	

	app.Listen(":3000")
}