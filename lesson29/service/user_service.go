package service

import (
	"context"
	"encoding/json"
	"lesson29/models"
	"lesson29/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo    *repository.UserRepo
	redisClient *redis.Client
}

func NewUserService(userRepo *repository.UserRepo, redisClient *redis.Client) *UserService {
	return &UserService{userRepo: userRepo, redisClient: redisClient}
}

func (s *UserService) GetUserById(ctx context.Context, userId string) (*models.User, string, error) {
	cachedData, err := s.redisClient.Get(ctx, userId).Result()
	if err == nil {
		user := &models.User{}
		if jsonErr := json.Unmarshal([]byte(cachedData), user); jsonErr != nil {
			return nil, "", jsonErr
		}
		return user, "cache", nil
	}

	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", nil
	}

	if jsonData, err := json.Marshal(user); err == nil {
		s.redisClient.Set(ctx, userId, jsonData, 10*time.Minute)
	}

	return user, "database", nil
}
