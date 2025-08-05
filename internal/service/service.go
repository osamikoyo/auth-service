package service

import (
	"auth/internal/entity/user"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	Repository interface{
		CreateUser(context.Context, *user.User) error
		UpdateUser(context.Context, string, string, interface{}) error
		GetUserByID(context.Context, string) (*user.User, error)
		GetUserByUsername(context.Context, string) (*user.User, error)
		DeleteUser(context.Context, string) error
	}

	Service struct{
		repository Repository

		timeout time.Duration
		key string
		timeForToken time.Duration
	}
)

var (
	ErrInternal = errors.New("internal server error")
	ErrInvalidInput = errors.New("invalid input")
	ErrIncorrectPassword = errors.New("incorrect password")
)

func NewService(repo Repository,key string, timeout, timeForToken time.Duration) *Service {
	return &Service{
		repository: repo,
		timeout: timeout,
		key: key,
		timeForToken: timeForToken,
	}
}

func (s *Service) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.timeout)
}

func (s *Service) Register(username, password string) (string, error) {
	if len(username) == 0 || len(password) == 0 {
		return "", ErrInvalidInput
	}
	
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", ErrInvalidInput
	}

	password = string(hashed)

	user := user.NewUser(username, password)

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repository.CreateUser(ctx, user);err != nil{
		return "", err
	}

	token, err := NewJwtKey(user.UID, s.key, s.timeForToken)
	if err != nil{
		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) Login(username, password string) (string, error) {
	if len(username) == 0 || len(password) == 0 {
		return "", ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil{
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));err != nil{
		return "", ErrIncorrectPassword
	}

	token, err := NewJwtKey(user.UID, s.key, s.timeForToken)
	if err != nil{
		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) GetAccount(uid string) (*user.User, error) {
	if len(uid) == 0{
		return nil, ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	user, err := s.repository.GetUserByID(ctx, uid)
	if err != nil{
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteAccount(uid string) error {
	if len(uid) == 0 {
		return ErrInvalidInput
	}

	ctx, cancel := s.context()
	defer cancel()

	if err := s.repository.DeleteUser(ctx, uid);err != nil{
		return err
	}

	return nil
}