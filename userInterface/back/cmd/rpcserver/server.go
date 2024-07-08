package rpcserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"log/slog"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/rpc"
	gorillaRpcJson "github.com/gorilla/rpc/json"
	"github.com/rosty-git/Smartchain-backend/internal/models"
	"github.com/rosty-git/Smartchain-backend/internal/usecase"
)

type RPCServer struct {
	port      string
	jwtSecret string
	useCase   usecase.UseCase
}

func MakeRPCServer(port string, jwtSecret string, useCase usecase.UseCase) (RPCServer, error) {
	return RPCServer{
		port:      port,
		jwtSecret: jwtSecret,
		useCase:   useCase,
	}, nil
}

type RPCService struct {
	useCase   usecase.UseCase
	jwtSecret string
}

func NewRPCService(jwtSecret string, useCase usecase.UseCase) *RPCService {
	return &RPCService{
		useCase:   useCase,
		jwtSecret: jwtSecret,
	}
}

func (rpcs RPCServer) Serve() error {
	s := rpc.NewServer()
	s.RegisterCodec(gorillaRpcJson.NewCodec(), "application/json")

	rpcService := NewRPCService(rpcs.jwtSecret, rpcs.useCase)

	err := s.RegisterService(rpcService, "")
	if err != nil {
		slog.Error("rpc register ", "err", err)
	}

	http.Handle("/rpc", s)

	l, err := net.Listen("tcp", ":"+rpcs.port)
	if err != nil {
		slog.Error("rpc listen ", "err", err)
		return err
	}

	server := &http.Server{
		Handler:      http.DefaultServeMux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = server.Serve(l)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http serve ", "err", err)
		return err
	}

	return nil
}

type RegistrationArgs struct {
	Email, Pass string
}

type RegistrationResult struct {
	ID string
}

func (s *RPCService) Registration(r *http.Request, args *RegistrationArgs, result *RegistrationResult) error {
	createdUser, err := s.useCase.CreateUser(args.Email, args.Pass)
	if err != nil {
		slog.Error("rpc register ", "err", err)

		return err
	}

	result.ID = createdUser.ID
	return nil
}

type LoginResult struct {
	Token string
}

func (s *RPCService) Login(r *http.Request, args *RegistrationArgs, result *LoginResult) error {
	token, err := s.useCase.GetToken(args.Email, args.Pass)
	if err != nil {
		slog.Error("rpc login ", "err", err)

		return err
	}

	result.Token = token
	return nil
}

func (s *RPCService) GetUserFromToken(tokenString string) (models.User, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(s.jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expiry date
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return models.User{}, errors.New("token expired")
		}

		userJson, ok := claims["user"]
		if !ok {
			slog.Error("AuthMiddleware")

			return models.User{}, errors.New("user field not found")
		}

		var user models.User
		err := json.Unmarshal([]byte(userJson.(string)), &user)
		if err != nil {
			slog.Error("AuthMiddleware", "error", err)

			return models.User{}, errors.New("json unmarshal error")
		}

		return user, nil
	} else {
		return models.User{}, errors.New("token invalid")
	}
}

type GetUserIamArgs struct {
	Token string
}

type GetUserIamResult struct {
	User models.User
}

func (s *RPCService) GetUserIam(r *http.Request, args *GetUserIamArgs, result *GetUserIamResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	result.User = user
	return nil
}

type TokenArgs struct {
	Token string
}

type CreateAccountResult struct {
	Address  string
	Password string
}

func (s *RPCService) CreateAccount(r *http.Request, args *TokenArgs, result *CreateAccountResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	address, password, err := s.useCase.CreateAccountForUser(user)
	if err != nil {
		return err
	}

	result.Address = address
	result.Password = password
	return nil
}

type GetAccountsResult struct {
	Accounts []models.Account
}

func (s *RPCService) GetAccounts(r *http.Request, args *TokenArgs, result *GetAccountsResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	accounts, err := s.useCase.GetAccountsByUserID(user.ID)
	if err != nil {
		return err
	}

	result.Accounts = accounts
	return nil
}

type DeleteAllAccountsResult struct {
	Success bool
}

func (s *RPCService) DeleteAllAccounts(r *http.Request, args *TokenArgs, result *DeleteAllAccountsResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	if !user.Admin {
		return errors.New("user is not admin")
	}

	err = s.useCase.DeleteAllAccounts()
	if err != nil {
		return err
	}

	result.Success = true
	return nil
}

type GetAccountBalanceArgs struct {
	Token     string
	AccountId string
}

type GetAccountBalanceResult struct {
	ExtId   string
	Balance string
}

func (s *RPCService) GetAccountBalance(r *http.Request, args *GetAccountBalanceArgs, result *GetAccountBalanceResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	extId, balance, err := s.useCase.GetAccountBalance(user, args.AccountId)
	if err != nil {
		return err
	}

	result.ExtId = extId
	result.Balance = balance.String()
	return nil
}

type FaucetArgs struct {
	Token    string
	Address  string
	Password string
}

type FaucetResult struct {
	Success bool
}

func (s *RPCService) Faucet(r *http.Request, args *FaucetArgs, result *FaucetResult) error {
	_, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	err = s.useCase.Faucet(args.Address, args.Password)
	if err != nil {
		return err
	}

	result.Success = true
	return nil
}

type GetFaucetContractAddressResult struct {
	Address string
}

func (s *RPCService) GetFaucetContractAddress(r *http.Request, args *TokenArgs, result *GetFaucetContractAddressResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	if !user.Admin {
		return errors.New("user is not admin")
	}

	address, err := s.useCase.GetParam("faucetContractAddress")
	if err != nil {
		return err
	}

	result.Address = address
	return nil
}

type SetFaucetContractAddressArgs struct {
	Token string
	Value string
}

type SetFaucetContractAddressResult struct {
	Value string
}

func (s *RPCService) SetFaucetContractAddress(r *http.Request, args *SetFaucetContractAddressArgs, result *SetFaucetContractAddressResult) error {
	user, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	if !user.Admin {
		return errors.New("user is not admin")
	}

	newValue, err := s.useCase.SetParam("faucetContractAddress", args.Value)
	if err != nil {
		return err
	}

	result.Value = newValue
	return nil
}

type GetBlocksResult struct {
	Blocks []usecase.Block
}

func (s *RPCService) GetBlocks(r *http.Request, args *TokenArgs, result *GetBlocksResult) error {
	_, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	blocks, err := s.useCase.GetBlocks()
	if err != nil {
		return err
	}

	result.Blocks = blocks
	return nil
}

type GetBlockByNumberResult struct {
	Block usecase.Block
}

func (s *RPCService) GetBlockByNumber(r *http.Request, args *TokenArgs, result *GetBlocksResult) error {
	_, err := s.GetUserFromToken(args.Token)
	if err != nil {
		return err
	}

	blocks, err := s.useCase.GetBlocks()
	if err != nil {
		return err
	}

	result.Blocks = blocks
	return nil
}
