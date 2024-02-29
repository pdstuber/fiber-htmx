package router

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/pdstuber/fiber-htmx/views"
)

type Contact struct {
	Name  string
	Email string
	ID    string
}

const shutdownTimeout = 5 * time.Second

type Server struct {
	listenPort string
	errChan    chan error
	fiberApp   *fiber.App
}

func New(listenPort string) *Server {
	engine := html.NewFileSystem(http.FS(views.Viewsfs), ".html")

	engine.AddFunc("dec", func(n int) int {
		return n - 1
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil, "layouts/main")
	})

	app.Use("/css", filesystem.New(filesystem.Config{
		Root:       http.FS(views.Viewsfs),
		PathPrefix: "css",
		Browse:     true,
	}))

	app.Get("/contacts", func(c *fiber.Ctx) error {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			return err
		}
		return c.Render("contacts", renderContacts(page))
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	return &Server{
		listenPort: listenPort,
		fiberApp:   app,
	}
}

func renderContacts(page int) fiber.Map {
	return fiber.Map{
		"NextPage": page + 1,
		"Contacts": []Contact{
			{
				Name:  "Agent Smith",
				Email: "void15@null.org",
				ID:    "123456",
			},
			{
				Name:  "Agent Smith",
				Email: "void16@null.org",
				ID:    "123457",
			},
			{
				Name:  "Agent Smith",
				Email: "void17@null.org",
				ID:    "123458",
			},
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.fiberApp.Listen(s.listenPort); err != nil {
			s.errChan <- err
		}
	}()

	select {
	case err := <-s.errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop() error {
	return s.fiberApp.ShutdownWithTimeout(shutdownTimeout)
}
