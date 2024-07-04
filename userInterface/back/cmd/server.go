package main

import (
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rosty-git/Smartchain-backend/internal/config"
	"github.com/rosty-git/Smartchain-backend/internal/controllers"
	"github.com/rosty-git/Smartchain-backend/internal/database"
	"github.com/rosty-git/Smartchain-backend/internal/middlewares"
	"github.com/rosty-git/Smartchain-backend/internal/repository"
	"github.com/rosty-git/Smartchain-backend/pkg/logger"
)

func main() {
	c := config.NewConfig()

	logger.InitLogger(c.GetEnv())

	db, closer, err := database.New(c.GetDsn(), false)
	if err != nil {
		slog.Error("Failed to connect to database")
	}
	defer closer()

	err = database.Initialize(db)
	if err != nil {
		slog.Error("Failed to initialize database")
	}

	// Initialize a new Fiber app
	app := fiber.New()

	jwt := middlewares.AuthMiddleware(c.GetJwtSecret())

	repo := repository.MakeRepository()

	controller := controllers.MakeController(repo, db, c)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Post("/api/v1/users/register", controller.Register)
	app.Post("/api/v1/users/login", controller.Login)
	app.Get("/api/v1/users/iam", jwt, controller.GetUserIam)
	app.Post("/api/v1/accounts", jwt, controller.CreateAccount)
	app.Get("/api/v1/accounts", jwt, controller.GetAccounts)
	app.Get("/api/v1/accounts/:accountId/balance", jwt, controller.GetAccountBalance)
	app.Post("/api/v1/faucet", jwt, controller.Faucet)

	// Start the server on port .env.FIBER_PORT
	log.Fatal(app.Listen(":" + c.GetFiberPort()))
}
