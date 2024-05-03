package db

import (
	"6_7/example/internals/app/models"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type UsersStorage struct {
	databasePool *pgxpool.Pool
}

func NewUsersStorage(pool *pgxpool.Pool) *UsersStorage {
	storage := new(UsersStorage)
	storage.databasePool = pool
	return storage

}

func (storage *UsersStorage) GetUsersList(ctx context.Context, nameFilter string) []models.User {
	query := "SELECT id, name, rank FROM users"
	args := make([]interface{}, 0)
	if nameFilter != "" {
		query += " WHERE name LIKE $1"
		args = append(args, fmt.Sprintf("%%%s%%", nameFilter))
	}

	var result []models.User

	err := pgxscan.Select(ctx, storage.databasePool, &result, query, args...)

	if err != nil {
		log.Errorln(err)
	}

	return result
}

func (storage *UsersStorage) GetUserById(ctx context.Context, id int64) models.User {
	query := "SELECT id, name, rank FROM users WHERE id = $1"

	var result models.User

	err := pgxscan.Get(ctx, storage.databasePool, &result, query, id)

	if err != nil {
		log.Errorln(err)
	}

	return result
}

func (storage *UsersStorage) CreateUser(ctx context.Context, user models.User) error {
	query := "INSERT INTO users(name, rank) VALUES ($1, $2)"

	_, err := storage.databasePool.Exec(ctx, query, user.Name, user.Rank) //транзакция не нужна, у нас только один запрос

	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}
