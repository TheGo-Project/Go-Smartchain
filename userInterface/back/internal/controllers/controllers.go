package controllers

import (
	"log/slog"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	"github.com/rosty-git/Smartchain-backend/internal/config"
	"github.com/rosty-git/Smartchain-backend/internal/models"
	"github.com/rosty-git/Smartchain-backend/internal/repository"
	"github.com/rosty-git/Smartchain-backend/internal/usecase"
	"gorm.io/gorm"
)

type Controller struct {
	repo      repository.Repository
	db        *gorm.DB
	jwtSecret string
	gethUrl   string
	useCase   usecase.UseCase
}

func MakeController(repo repository.Repository, db *gorm.DB, c *config.Config, useCase usecase.UseCase) Controller {
	return Controller{
		repo:      repo,
		db:        db,
		jwtSecret: c.GetJwtSecret(),
		gethUrl:   c.GetGethUrl(),
		useCase:   useCase,
	}
}

type UserEmailPasswordForm struct {
	Email string `json:"email" form:"email"`
	Pass  string `json:"password" form:"password"`
}

func (ctrl Controller) Register(c *fiber.Ctx) error {
	slog.Info("Register")

	uepf := UserEmailPasswordForm{}

	if err := c.BodyParser(&uepf); err != nil {
		return err
	}

	createdUser, err := ctrl.useCase.CreateUser(uepf.Email, uepf.Pass)
	if err != nil {
		return err
	}

	slog.Info("Created user", "user", createdUser)

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func (ctrl Controller) Login(c *fiber.Ctx) error {
	slog.Info("Login")

	u := UserEmailPasswordForm{}

	if err := c.BodyParser(&u); err != nil {
		return err
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := ctrl.useCase.GetToken(u.Email, u.Pass)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func (ctrl Controller) GetUserIam(c *fiber.Ctx) error {
	user := c.Locals("User")

	return c.JSON(user)
}

func (ctrl Controller) GetAccounts(c *fiber.Ctx) error {
	slog.Info("GetAccounts")

	user := c.Locals("User").(models.User)

	accounts, err := ctrl.useCase.GetAccountsByUserID(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(accounts)
}

func (ctrl Controller) DeleteAllAccounts(c *fiber.Ctx) error {
	err := ctrl.useCase.DeleteAllAccounts()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{})
}

func (ctrl Controller) GetAccountBalance(c *fiber.Ctx) error {
	slog.Info("GetAccountBalance")

	user := c.Locals("User").(models.User)

	accountExtId, balance, err := ctrl.useCase.GetAccountBalance(user, c.Params("accountId"))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"ext_id":  accountExtId,
		"balance": balance.String(),
	})
}

func (ctrl Controller) CreateAccount(c *fiber.Ctx) error {
	slog.Info("CreateAccount")

	user := c.Locals("User").(models.User)

	slog.Info("CreateAccount", "user", user.ID)

	address, password, err := ctrl.useCase.CreateAccountForUser(user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"address":  address,
		"password": password,
	})
}

func (ctrl Controller) Faucet(c *fiber.Ctx) error {
	slog.Info("Faucet")

	type FaucetRequest struct {
		Address  string `json:"address" form:"address"`
		Password string `json:"password" form:"password"`
	}

	fr := FaucetRequest{}

	if err := c.BodyParser(&fr); err != nil {
		return err
	}

	err := ctrl.useCase.Faucet(fr.Address, fr.Password)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{})
}

func (ctrl Controller) GetFaucetContractAddress(c *fiber.Ctx) error {
	param, err := ctrl.useCase.GetParam("faucetContractAddress")
	if err != nil {
		slog.Error("Failed to get faucet contract address", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"value": param})
}

func (ctrl Controller) SetParam(c *fiber.Ctx, key string) (string, error) {
	type SetRequest struct {
		Value string `json:"value" form:"value"`
	}

	sfr := SetRequest{}

	if err := c.BodyParser(&sfr); err != nil {
		slog.Error("Failed to get faucet contract address", "err", err)

		return "", err
	}

	_, err := ctrl.repo.SetParam(ctrl.db, key, sfr.Value)
	if err != nil {
		slog.Error("Failed to set faucet contract address", "err", err)

		return "", err
	}

	return sfr.Value, nil
}

func (ctrl Controller) SetFaucetContractAddress(c *fiber.Ctx) error {
	type SetRequest struct {
		Value string `json:"value" form:"value"`
	}

	sfr := SetRequest{}

	value, err := ctrl.useCase.SetParam("faucetContractAddress", sfr.Value)
	if err != nil {
		slog.Error("Failed to set faucet contract address", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"value": value})
}

//func (ctrl Controller) GetMainAccountAddress(c *fiber.Ctx) error {
//	param, err := ctrl.GetParam("mainAccountAddress")
//	if err != nil {
//		slog.Error("Failed to get faucet contract address", "err", err)
//
//		return c.Status(500).SendString(err.Error())
//	}
//
//	return c.JSON(fiber.Map{"value": param})
//}
//
//func (ctrl Controller) SetMainAccountAddress(c *fiber.Ctx) error {
//	value, err := ctrl.SetParam(c, "mainAccountAddress")
//	if err != nil {
//		slog.Error("Failed to set faucet contract address", "err", err)
//
//		return c.Status(500).SendString(err.Error())
//	}
//
//	return c.JSON(fiber.Map{"value": value})
//}

func (ctrl Controller) HealthCheck(c *fiber.Ctx) error {
	count, err := ctrl.repo.GetUsersCount(ctrl.db)
	if err != nil {
		slog.Error("Failed to set faucet contract address", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"count": count, "healthy": true})
}

func (ctrl Controller) GetBlocks(c *fiber.Ctx) error {
	blocks, err := ctrl.useCase.GetBlocks()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"blocks": blocks})
}

func (ctrl Controller) GetBlockByNumber(c *fiber.Ctx) error {
	client, err := ethclient.Dial(ctrl.gethUrl)
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)
	}
	defer client.Close()

	i, err := strconv.ParseInt(c.Params("blockNumber"), 10, 64)
	if err != nil {
		slog.Error("Failed to parse block number", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	block, err := ctrl.useCase.GetBlockByNumber(big.NewInt(i))
	if err != nil {
		slog.Error("Failed to retrieve block %d", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(block)
}
