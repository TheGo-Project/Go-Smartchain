package controllers

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rosty-git/Smartchain-backend/internal/config"
	"github.com/rosty-git/Smartchain-backend/internal/models"
	"github.com/rosty-git/Smartchain-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controller struct {
	repo      repository.Repository
	db        *gorm.DB
	jwtSecret string
	gethUrl   string
}

func MakeController(repo repository.Repository, db *gorm.DB, c *config.Config) Controller {
	return Controller{
		repo:      repo,
		db:        db,
		jwtSecret: c.GetJwtSecret(),
		gethUrl:   c.GetGethUrl(),
	}
}

type User struct {
	Email string `json:"email" form:"email"`
	Pass  string `json:"password" form:"password"`
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (ctrl Controller) Register(c *fiber.Ctx) error {
	slog.Info("Register")

	u := User{}

	if err := c.BodyParser(&u); err != nil {
		return err
	}

	user := models.User{
		Email: u.Email,
	}

	// Generate salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	user.Salt = hex.EncodeToString(salt)
	slog.Info("Registering user", "user.Salt", user.Salt)

	// Hash password with salt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Pass+user.Salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	usersCount, err := ctrl.repo.GetUsersCount(ctrl.db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if usersCount == 0 {
		user.Admin = true
	}

	createdUser, err := ctrl.repo.CreateUser(ctrl.db, user)
	if err != nil {
		return err
	}

	slog.Info("Created user", "user", createdUser)

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func (ctrl Controller) Login(c *fiber.Ctx) error {
	slog.Info("Register")

	jwtTtl := time.Duration(604800) * 1_000_000_000
	//cookieName := "Authorization"

	u := User{}

	if err := c.BodyParser(&u); err != nil {
		return err
	}

	user, err := ctrl.repo.GetUserByEmail(ctrl.db, u.Email)
	if err != nil {
		return err
	}

	// Compare hashed passwords
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Pass+user.Salt)); err != nil {
		slog.Error("Wrong password", "err", err)

		return err
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(jwtTtl).Unix(),
		"user": string(userJson),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(ctrl.jwtSecret))
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

	accounts, err := ctrl.repo.GetAccountsByUserID(ctrl.db, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(accounts)
}

func (ctrl Controller) GetAccountBalance(c *fiber.Ctx) error {
	slog.Info("GetAccountBalance")

	user := c.Locals("User").(models.User)

	account, err := ctrl.repo.GetAccountByID(ctrl.db, c.Params("accountId"))
	if err != nil {
		return err
	}

	if account.UserID != user.ID {
		return c.JSON(fiber.Map{})
	}

	slog.Info("GetAccountBalance", "Dial", ctrl.gethUrl)

	client, err := ethclient.Dial(ctrl.gethUrl)
	if err != nil {
		slog.Info("Failed to connect to the Ethereum client", "err", err)

		return c.Status(500).SendString(err.Error())
	}
	defer client.Close()

	accountAddress := common.HexToAddress(account.ExtID)

	balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		slog.Info("Failed to get balance", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	slog.Info("GetAccountBalance", "balance", balance)

	return c.JSON(fiber.Map{
		"ext_id":  account.ExtID,
		"balance": balance.String(),
	})
}

func (ctrl Controller) CreateAccount(c *fiber.Ctx) error {
	slog.Info("CreateAccount")

	user := c.Locals("User").(models.User)

	slog.Info("CreateAccount", "user", user.ID)

	client, err := ethclient.Dial(ctrl.gethUrl)
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)
	}
	defer client.Close()

	slog.Info("createAccount", "client", client.Client())
	id, err := client.ChainID(context.TODO())
	if err != nil {
		slog.Info("Failed to get the chain ID", "err", err)

		return c.Status(500).SendString(err.Error())
	}
	slog.Info("CreateAccount", "ChainID", id)

	slog.Info("Connected to Ethereum client", "client", client)

	password := generateRandomString(30)

	slog.Info("CreateAccount", "password", password)

	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		slog.Error("Failed to create new account", "err", err)
	}

	fmt.Printf("New account created: %s\n", account.Address.Hex())

	keystorePath := filepath.Join(strings.Replace(account.URL.String(), "keystore://", "", 1))
	keystoreBytes, err := os.ReadFile(keystorePath)
	if err != nil {
		slog.Error("Failed to read keystore file", "err", err)
	}

	e := os.Remove(keystorePath)
	if e != nil {
		slog.Error("Failed to remove keystore file", "err", e)
	}

	var keystoreData models.KeyStore
	err = json.Unmarshal(keystoreBytes, &keystoreData)
	if err != nil {
		slog.Error("Failed to unmarshal keystore JSON", "err", err)
	}

	newAccount := models.Account{
		UserID:   user.ID,
		ExtID:    account.Address.Hex(),
		Keystore: keystoreData,
	}

	result := ctrl.db.Create(&newAccount)
	if result.Error != nil {
		slog.Error("Failed to create account", "err", result.Error)
	}

	return c.JSON(fiber.Map{
		"address":  account.Address.Hex(),
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

	const faucetABI = `[
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "Deposit",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "Withdrawal",
			"type": "event"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "withdraw_amount",
					"type": "uint256"
				}
			],
			"name": "withdraw",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"stateMutability": "payable",
			"type": "receive"
		}
	]`

	client, err := ethclient.Dial(ctrl.gethUrl)
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)
	}
	defer client.Close()

	contractAbi, err := abi.JSON(strings.NewReader(faucetABI))
	if err != nil {
		slog.Error("Failed to parse contract ABI", "err", err)
	}

	faucetContractAddress, err := ctrl.repo.GetParam(ctrl.db, "faucetContractAddress")
	if err != nil || faucetContractAddress == "" {
		slog.Info("faucetContractAddress", "err", err)

		return c.Status(500).SendString(err.Error())
	}

	address := common.HexToAddress(faucetContractAddress)
	contract := bind.NewBoundContract(address, contractAbi, client, client, client)

	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		slog.Error("Failed to retrieve contract balance", "err", err)
	}
	slog.Info("Faucet", "Contract balance", balance.String())

	account, err := ctrl.repo.GetAccountsByExtID(ctrl.db, fr.Address)
	if err != nil {
		slog.Error("Failed to retrieve account", "err", err)
	}

	accountKeystoreBytes, err := json.Marshal(account.Keystore)
	if err != nil {
		slog.Error("Failed to marshal keystore JSON", "err", err)
	}

	privateKey, err := keystore.DecryptKey(accountKeystoreBytes, fr.Password)
	if err != nil {
		slog.Error("Failed to decrypt private key", "err", err)
	}

	amount := big.NewInt(0.5e18) // 0.5 ETH в wei
	withdraw(client, privateKey.PrivateKey, contract, amount)

	return c.JSON(fiber.Map{})
}

func withdraw(client *ethclient.Client, privateKey *ecdsa.PrivateKey, contract *bind.BoundContract, amount *big.Int) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		slog.Error("Failed to create authorized transactor", "err", err)
	}

	tx, err := contract.Transact(auth, "withdraw", amount)
	if err != nil {
		slog.Error("Failed to request withdrawal", "err", err)
	}

	fmt.Printf("Withdrawal transaction sent: %s\n", tx.Hash().Hex())
}

func (ctrl Controller) GetParam(key string) (string, error) {
	param, err := ctrl.repo.GetParam(ctrl.db, key)
	if err != nil {
		slog.Error("Failed to get faucet contract address", "err", err)

		return "", err
	}
	return param, nil
}

func (ctrl Controller) GetFaucetContractAddress(c *fiber.Ctx) error {
	param, err := ctrl.GetParam("faucetContractAddress")
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
	value, err := ctrl.SetParam(c, "faucetContractAddress")
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
