package main

import (
	"context"
	"fmt"
	"music-master/config"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	musictrackcustomer "music-master/internal/api/v1/customer/musictrack"
	playlistcustomer "music-master/internal/api/v1/customer/playlist"
	"music-master/internal/db"
	"music-master/internal/db/elasticsearch"
	"music-master/internal/util/converter"
	"music-master/internal/util/server"
	_ "music-master/internal/util/swagger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Connect to mongo
	mongoDB, err := db.New(cfg)
	if err != nil {
		panic(err)
	}
	defer mongoDB.Disconnect()

	es, err := elasticsearch.NewESClient(cfg)
	if err != nil {
		panic(err)
	}

	musicTrackCollection := db.NewMusicTrackCollection(mongoDB)
	playlistCollection := db.NewPlaylistCollection(mongoDB)

	musicTrackES := elasticsearch.NewMusicTrackCollection(es)

	fmt.Println("cfg", cfg)
	// Init HTTP server
	e := New(&Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	converter := converter.NewModelConverter()
	musicTrackCustomer := musictrackcustomer.New(musicTrackCollection, converter, musicTrackES)
	playlistCustomer := playlistcustomer.New(playlistCollection, converter)

	v1cRouter := e.Group("/v1")

	// * customer
	v1cRouter = v1cRouter.Group("/customer")
	musictrackcustomer.NewHTTP(musicTrackCustomer, nil, v1cRouter.Group("/music-tracks"))
	playlistcustomer.NewHTTP(playlistCustomer, nil, v1cRouter.Group("/playlists"))

	// Static page for Swagger API specs
	if cfg.Debug {
		e.Static("/swaggerui", "swaggerui")
	}

	Start(e)
}

// Config represents server specific config
type Config struct {
	Stage        string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	Debug        bool
	AllowOrigins []string
}

// New instances new Echo server
func New(cfg *Config) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{LogLevel: log.ERROR}),
		Headers(), CORS([]string{"*"}))
	e.GET("/", healthCheck)
	e.Validator = server.NewValidator()
	e.HTTPErrorHandler = server.NewErrorHandler(e).Handle
	e.Debug = cfg.Debug

	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

// Start starts echo server
func Start(e *echo.Echo) {
	// graceful shutdown for dev environment

	// Start server
	go func() {
		if err := e.StartServer(e.Server); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("shutting down the server")
			} else {
				fmt.Println("error shutting down the server: ", err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
	})
}

// CORS adds Cross-Origin Resource Sharing support
func CORS(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400,
	})
}
