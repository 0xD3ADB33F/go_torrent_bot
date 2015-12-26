package main

import (
	"fmt"
	"gopkg.in/redis.v3"
)

type torrentsRepository struct {
	prefix string
	client *redis.Client
}

func newTorrentsRepository(prefix string, client *redis.Client) *torrentsRepository {
	return &torrentsRepository{prefix, client}
}

func (repo *torrentsRepository) Add(chatID, userID int, content []byte) {
	repo.client.Set(repo.key(chatID, userID), content, 0)
}

func (repo *torrentsRepository) Exists(chatID, userID int) (result bool) {
	result, _ = repo.client.Exists(repo.key(chatID, userID)).Result()
	return
}

func (repo *torrentsRepository) Get(chatID, userID int) ([]byte, error) {
	return repo.client.Get(repo.key(chatID, userID)).Bytes()
}

func (repo *torrentsRepository) Delete(chatID, userID int) {
	repo.client.Del(repo.key(chatID, userID))
}

func (repo *torrentsRepository) key(chatID, userID int) string {
	return fmt.Sprintf("%s_%d_%d", repo.prefix, chatID, userID)
}
