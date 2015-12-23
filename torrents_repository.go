package main

import (
	"fmt"
	"gopkg.in/redis.v3"
)

type TorrentsRepository struct {
	prefix string
	client *redis.Client
}

func newTorrentsRepository(prefix string, client *redis.Client) *TorrentsRepository {
	return &TorrentsRepository{prefix, client}
}

func (repo *TorrentsRepository) Add(chatID, userID int, content []byte) {
	repo.client.Set(repo.key(chatID, userID), content, 0)
}

func (repo *TorrentsRepository) Exists(chatID, userID int) (result bool) {
	result, _ = repo.client.Exists(repo.key(chatID, userID)).Result()
	return
}

func (repo *TorrentsRepository) Get(chatID, userID int) ([]byte, error) {
	return repo.client.Get(repo.key(chatID, userID)).Bytes()
}

func (repo *TorrentsRepository) Delete(chatID, userID int) {
	repo.client.Del(repo.key(chatID, userID))
}

func (repo *TorrentsRepository) key(chatID, userID int) string {
	return fmt.Sprintf("%s_%d_%d", repo.prefix, chatID, userID)
}
