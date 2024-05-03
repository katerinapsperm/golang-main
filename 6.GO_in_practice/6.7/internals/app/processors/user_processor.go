package processors

import (
	"6_7/example/internals/app/db"
	"6_7/example/internals/app/models"
	"context"
	"errors"
)

type UsersProcessor struct {
	storage *db.UsersStorage
}

func NewUsersProcessor(storage *db.UsersStorage) *UsersProcessor {
	processor := new(UsersProcessor)
	processor.storage = storage
	return processor
}

func (processor *UsersProcessor) CreateUser(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.New("name should not be empty")
	}

	return processor.storage.CreateUser(ctx, user)
}

func (processor *UsersProcessor) FindUser(ctx context.Context, id int64) (models.User, error) {
	user := processor.storage.GetUserById(ctx, id)

	if user.Id != id {
		return user, errors.New("user not found")
	}

	return user, nil

}

func (processor *UsersProcessor) ListUsers(ctx context.Context, nameFilter string) ([]models.User, error) {
	return processor.storage.GetUsersList(ctx, nameFilter), nil
}
