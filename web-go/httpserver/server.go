package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/martin-flower/roboz-web-go/config"
	"github.com/martin-flower/roboz-web-go/handlers"
	"github.com/spf13/viper"
)

type Server struct {
	app *fiber.App
}

func New(timeouts config.Timeouts) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  timeouts.ReadTimeout,
		WriteTimeout: timeouts.WriteTimeout,
		IdleTimeout:  timeouts.IdleTimeout,
	})
	return &Server{
		app: app,
	}
}

func (srv Server) Setup() {
	srv.app.Get("/health", handlers.Health)
	srv.app.Get("/list", handlers.List)
	srv.app.Post("/developer-test/enter-path", handlers.Enter)

	// documentation: default landing page
	srv.app.Get("/", redirectToSwagger)
	srv.app.Get("/docs/swagger/*", swagger.HandlerDefault)
}

func redirectToSwagger(c *fiber.Ctx) error {
	return c.Redirect("/docs/swagger/index.html")
}

func (srv Server) Shutdown() error {
	return srv.app.Shutdown()
}

func (srv Server) Listen() error {
	return srv.app.Listen(":" + viper.GetString("port"))
}
