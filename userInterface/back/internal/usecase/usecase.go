package usecase

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"github.com/golang-jwt/jwt/v4"
	"github.com/rosty-git/Smartchain-backend/internal/config"
	"github.com/rosty-git/Smartchain-backend/internal/models"
	"github.com/rosty-git/Smartchain-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UseCase struct {
	repo   repository.Repository
	db     *gorm.DB
	config *config.Config
}

func MakeUseCase(repo repository.Repository, db *gorm.DB, c *config.Config) UseCase {
	return UseCase{
		repo:   repo,
		db:     db,
		config: c,
	}
}

func (uc UseCase) CreateUser(email, password string) (models.User, error) {
	user := models.User{
		Email: email,
	}

	// Generate salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return models.User{}, err
	}
	user.Salt = hex.EncodeToString(salt)

	// Hash password with salt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password+user.Salt), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user.Password = string(hashedPassword)

	usersCount, err := uc.repo.GetUsersCount(uc.db)
	if err != nil {
		return models.User{}, err
	}

	if usersCount == 0 {
		user.Admin = true
	}

	createdUser, err := uc.repo.CreateUser(uc.db, user)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (uc UseCase) GetToken(email, password string) (string, error) {
	jwtTtl := time.Duration(604800) * 1_000_000_000

	user, err := uc.repo.GetUserByEmail(uc.db, email)
	if err != nil {
		return "", err
	}

	// Compare hashed passwords
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt)); err != nil {
		slog.Error("Wrong password", "err", err)

		return "", err
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(jwtTtl).Unix(),
		"user": string(userJson),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(uc.config.GetJwtSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc UseCase) GetAccountsByUserID(id string) ([]models.Account, error) {
	accounts, err := uc.repo.GetAccountsByUserID(uc.db, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (uc UseCase) DeleteAllAccounts() error {
	err := uc.repo.DeleteAllAccounts(uc.db)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (uc UseCase) CreateAccountForUser(user models.User) (string, string, error) {
	client, err := ethclient.Dial(uc.config.GetGethUrl())
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)
	}
	defer client.Close()

	slog.Info("createAccount", "client", client.Client())
	id, err := client.ChainID(context.TODO())
	if err != nil {
		slog.Info("Failed to get the chain ID", "err", err)

		return "", "", err
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

	result := uc.db.Create(&newAccount)
	if result.Error != nil {
		slog.Error("Failed to create account", "err", result.Error)
	}

	return account.Address.Hex(), password, nil
}

func (uc UseCase) GetAccountBalance(user models.User, accountId string) (string, *big.Int, error) {
	account, err := uc.repo.GetAccountByID(uc.db, accountId)
	if err != nil {
		return "", nil, err
	}

	if account.UserID != user.ID {
		return "", nil, errors.New("")
	}

	slog.Info("GetAccountBalance", "Dial", uc.config.GetGethUrl())

	client, err := ethclient.Dial(uc.config.GetGethUrl())
	if err != nil {
		slog.Info("Failed to connect to the Ethereum client", "err", err)

		return "", nil, err
	}
	defer client.Close()

	accountAddress := common.HexToAddress(account.ExtID)

	balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		slog.Info("Failed to get balance", "err", err)

		return "", nil, err
	}

	slog.Info("GetAccountBalance", "balance", balance)

	return account.UserID, balance, nil
}

func withdraw(client *ethclient.Client, privateKey *ecdsa.PrivateKey, contract *bind.BoundContract, amount *big.Int) error {
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	slog.Error("Failed to get suggested gas price", "err", err)
	//	return err
	//}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		slog.Error("Failed to create authorized transactor", "err", err)
	}

	//gasLimit := uint64(200000)
	//
	//auth.Nonce = nil
	//auth.GasLimit = gasLimit
	//auth.GasPrice = gasPrice

	tx, err := contract.Transact(auth, "withdraw", amount)
	if err != nil {
		slog.Error("Failed to request withdrawal", "err", err)

		return err
	}

	fmt.Printf("Withdrawal transaction sent: %s\n", tx.Hash().Hex())

	return nil
}

func (uc UseCase) Faucet(targetAccountAddress string, password string) error {
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

	client, err := ethclient.Dial(uc.config.GetGethUrl())
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)
	}
	defer client.Close()

	contractAbi, err := abi.JSON(strings.NewReader(faucetABI))
	if err != nil {
		slog.Error("Failed to parse contract ABI", "err", err)
	}

	faucetContractAddress, err := uc.repo.GetParam(uc.db, "faucetContractAddress")
	if err != nil || faucetContractAddress == "" {
		slog.Info("faucetContractAddress", "err", err)

		return err
	}

	addressHex := common.HexToAddress(faucetContractAddress)
	contract := bind.NewBoundContract(addressHex, contractAbi, client, client, client)

	faucetBalance, err := client.BalanceAt(context.Background(), addressHex, nil)
	if err != nil {
		slog.Error("Failed to retrieve contract balance", "err", err)
	}
	slog.Info("Faucet", "Contract balance", faucetBalance.String())

	accountAddress := common.HexToAddress(targetAccountAddress)
	accountBalance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		slog.Error("Failed to retrieve contract balance", "err", err)
	}
	slog.Info("Faucet", "Account balance", accountBalance.String())

	account, err := uc.repo.GetAccountsByExtID(uc.db, targetAccountAddress)
	if err != nil {
		slog.Error("Failed to retrieve account", "err", err)
	}

	accountKeystoreBytes, err := json.Marshal(account.Keystore)
	if err != nil {
		slog.Error("Failed to marshal keystore JSON", "err", err)
	}

	privateKey, err := keystore.DecryptKey(accountKeystoreBytes, password)
	if err != nil {
		slog.Error("Failed to decrypt private key", "err", err)
	}

	amount := big.NewInt(0.5e18) // 0.5 ETH в wei
	err = withdraw(client, privateKey.PrivateKey, contract, amount)
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) GetParam(key string) (string, error) {
	param, err := uc.repo.GetParam(uc.db, key)
	if err != nil {
		slog.Error("Failed to get faucet contract address", "err", err)

		return "", err
	}
	return param, nil
}

func (uc UseCase) SetParam(key string, value string) (string, error) {
	_, err := uc.repo.SetParam(uc.db, key, value)
	if err != nil {
		slog.Error("Failed to set faucet contract address", "err", err)

		return "", err
	}

	return value, nil
}

type Block struct {
	Number            uint64 `json:"number"`
	Hash              string `json:"hash"`
	ParentHash        string `json:"parentHash"`
	Time              uint64 `json:"time"`
	TransactionsCount int    `json:"transactionsCount"`
}

func (uc UseCase) GetBlocks() ([]Block, error) {
	client, err := ethclient.Dial(uc.config.GetGethUrl())
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)

		return nil, err
	}
	defer client.Close()

	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		slog.Error("Failed to retrieve latest block number", "err", err)

		return nil, err
	}

	var blocks []Block

	var n uint64 = 20

	if latestBlockNumber < n {
		n = latestBlockNumber
	}

	for i := latestBlockNumber; i > latestBlockNumber-n; i-- {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			slog.Error("Failed to retrieve block %d", "err", err)
		}

		blocks = append(blocks, Block{
			Number: block.Number().Uint64(),
			Hash:   block.Hash().Hex(),
			//ParentHash: block.ParentHash().Hex(),
			Time: block.Time(),
		})
	}

	return blocks, nil
}

func (uc UseCase) GetBlockByNumber(blockNumber *big.Int) (Block, error) {
	client, err := ethclient.Dial(uc.config.GetGethUrl())
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "err", err)

		return Block{}, err
	}
	defer client.Close()

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		slog.Error("Failed to retrieve block %d", "err", err)

		return Block{}, err
	}

	return Block{
		Number:            block.Number().Uint64(),
		Hash:              block.Hash().Hex(),
		ParentHash:        block.ParentHash().Hex(),
		Time:              block.Time(),
		TransactionsCount: len(block.Transactions()),
	}, nil
}
