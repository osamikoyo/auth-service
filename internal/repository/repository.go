package repository

import (
	"auth/internal/entity/user"
	"auth/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrAlreadyExist     = errors.New("user already exist")
	ErrRepositoryFailed = errors.New("internal error")
	ErrNotFound         = errors.New("user not found")
)

type Repository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, usr *user.User) error {
	var existingUser user.User
	if err := r.db.Where("username = ?", usr.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user with username %s already exists", usr.Username)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	if err := r.db.Create(usr).Error; err != nil {
		r.logger.Error("failed create user", zap.Error(err))
		return fmt.Errorf("error creating user: %v", err)
	}

	r.logger.Info("user created", zap.Any("user", usr))

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, uid string, key string, update interface{}) error {
	if rows, err := gorm.G[user.User](r.db).Where("uid = ?", uid).Update(ctx, key, update); err != nil {
		r.logger.Error("failed update user",
			zap.String("key", key),
			zap.Any("update", update),
			zap.Error(err))

		if rows == 0 {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	r.logger.Info("user updated",
		zap.String("uid", uid),
		zap.String("key", key),
		zap.Any("update", update))

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, uid string) (*user.User, error) {
	user, err := gorm.G[user.User](r.db).Where("uid = ?", uid).First(ctx)
	if err != nil {
		r.logger.Error("failed fetch user",
			zap.String("uid", uid),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrRepositoryFailed
	}

	r.logger.Info("user fetched", zap.Any("user", user))

	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	user, err := gorm.G[user.User](r.db).Where("username = ?", username).First(ctx)
	if err != nil {
		r.logger.Error("failed fetch user",
			zap.String("username", username),
			zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrRepositoryFailed
	}

	r.logger.Info("user fetched", zap.Any("user", user))

	return &user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, uid string) error {
	if rows, err := gorm.G[user.User](r.db).Where("uid = uid", uid).Delete(ctx); err != nil {
		r.logger.Error("failed delete user",
			zap.String("uid", uid),
			zap.Error(err))

		if rows == 0 {
			return ErrNotFound
		}

		return ErrRepositoryFailed
	}

	r.logger.Info("user deleted", zap.String("uid", uid))

	return nil
}
