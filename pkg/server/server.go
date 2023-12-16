package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/fatfatcocofat/rosamsoe/platform/logger"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

type Server struct {
	App    *fiber.App
	Config *config.Config
	Logger *zerolog.Logger
}

func New(config *config.Config) *Server {
	srv := &Server{
		App: fiber.New(fiber.Config{
			ServerHeader: "Rosamsoe",
		}),
		Config: config,
		Logger: &logger.Logger,
	}

	srv.App.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger.Logger,
	}))

	srv.App.Use(cors.New())

	srv.routes()
	srv.handleSwaggerRoutes()
	srv.catchNotFoundRoutes()

	return srv
}

func (s *Server) Serve() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigch
		logger.Info().Msg("Shutting down server....")
		_ = s.App.ShutdownWithTimeout(60 * time.Second)
	}()

	listenAddr := fmt.Sprintf("%s:%d", s.Config.ServerHost, s.Config.ServerPort)
	if err := s.App.Listen(listenAddr); err != nil {
		logger.Fatal().Err(err).Msg("Oops... server is not running")
	}
}
