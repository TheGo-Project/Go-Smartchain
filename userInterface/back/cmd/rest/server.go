package main

import (
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rosty-git/Smartchain-backend/cmd/rpcserver"
	"github.com/rosty-git/Smartchain-backend/internal/config"
	"github.com/rosty-git/Smartchain-backend/internal/controllers"
	"github.com/rosty-git/Smartchain-backend/internal/database"
	"github.com/rosty-git/Smartchain-backend/internal/middlewares"
	"github.com/rosty-git/Smartchain-backend/internal/repository"
	"github.com/rosty-git/Smartchain-backend/internal/usecase"
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

	jwt := middlewares.AuthMiddleware(c.GetJwtSecret())
	adm := middlewares.AdmMiddleware()

	repo := repository.MakeRepository()
	useCase := usecase.MakeUseCase(repo, db, c)
	controller := controllers.MakeController(repo, db, c, useCase)

	rpcServer, err := rpcserver.MakeRPCServer(c.GetRpcPort(), c.GetJwtSecret(), useCase)
	if err != nil {
		slog.Error("Failed to initialize RPC server")
	}

	go func() {
		err := rpcServer.Serve()
		if err != nil {
			slog.Error("Failed to start RPC server")
		}
	}()

	// Initialize a new Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000, http://34.122.25.84",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Post("/api/v1/users/register", controller.Register)
	app.Post("/api/v1/users/login", controller.Login)
	app.Get("/api/v1/health/check", controller.HealthCheck)
	app.Get("/api/v1/users/iam", jwt, controller.GetUserIam)
	app.Post("/api/v1/accounts", jwt, controller.CreateAccount)
	app.Get("/api/v1/accounts", jwt, controller.GetAccounts)
	app.Delete("/api/v1/accounts/all", jwt, adm, controller.DeleteAllAccounts)
	app.Get("/api/v1/accounts/:accountId/balance", jwt, controller.GetAccountBalance)
	app.Post("/api/v1/faucet", jwt, controller.Faucet)
	app.Get("/api/v1/params/faucet-contract-address", jwt, adm, controller.GetFaucetContractAddress)
	app.Post("/api/v1/params/faucet-contract-address", jwt, adm, controller.SetFaucetContractAddress)
	//app.Get("/api/v1/params/main-account-address", jwt, adm, controller.GetMainAccountAddress)
	//app.Post("/api/v1/params/main-account-address", jwt, adm, controller.SetMainAccountAddress)
	app.Get("/api/v1/blocks", jwt, controller.GetBlocks)
	app.Get("/api/v1/blocks/:blockNumber", jwt, controller.GetBlockByNumber)

	// Start the server on port .env.FIBER_PORT
	log.Fatal(app.Listen(":" + c.GetFiberPort()))
}
