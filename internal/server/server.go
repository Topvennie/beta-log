// Package server starts the server
package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/Topvennie/beta-log/internal/database/repository"
	"github.com/Topvennie/beta-log/internal/server/api"
	"github.com/Topvennie/beta-log/internal/server/dto"
	middlewares "github.com/Topvennie/beta-log/internal/server/middlewares"
	"github.com/Topvennie/beta-log/pkg/config"
	zapfiber "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/shareed2k/goth_fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New() (*Server, error) {
	pool := repository.Pool()

	// Construct app
	app := fiber.New(fiber.Config{
		BodyLimit:         20 * 1024 * 1024,
		ReadBufferSize:    8096,
		StreamRequestBody: true,
		ErrorHandler:      middlewares.ErrorHandler(),
	})

	app.Use(recover.New())
	app.Use(middlewares.BodyCapture)
	app.Use(zapfiber.New(zapfiber.Config{
		Logger: zap.L(),
	}))
	if config.IsDev() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
		}))
	}

	// Session storage
	sessionStore := postgres.New(postgres.Config{
		DB: pool,
	})

	sessionCookieName := config.GetString("app.name")
	if sessionCookieName == "" {
		sessionCookieName = "session"
	}

	handler, store := session.NewWithStore(session.Config{
		CookieSecure:    !config.IsDev(),
		CookieSameSite:  "Lax",
		CookieHTTPOnly:  true,
		Storage:         sessionStore,
		IdleTimeout:     24 * time.Hour,
		AbsoluteTimeout: 7 * 24 * time.Hour,
		Extractor:       extractors.FromCookie(sessionCookieName + "_session"),
	})
	goth_fiber.SessionManager = goth_fiber.NewSessionManager(store)

	// Init dto validator
	if err := dto.InitValidator(); err != nil {
		return nil, err
	}

	// API
	if err := api.New(app.Group("/api").Use(handler)); err != nil {
		return nil, fmt.Errorf("initialize api %w", err)
	}

	// Static files if served in production
	if !config.IsDev() {
		app.Get("/*", static.New("./public"))
		app.Get("*", func(c fiber.Ctx) error {
			return c.SendFile("./public/index.html")
		})
	}

	// Fallback
	app.All("/api/*", func(c fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := config.GetDefaultInt("server.port", 3000)
	host := config.GetDefaultString("server.host", "0.0.0.0")

	srv := &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}

	return srv, nil
}
