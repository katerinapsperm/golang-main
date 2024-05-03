package db

import (
	"6_7/example/internals/app/models"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type CarsStorage struct {
	databasePool *pgxpool.Pool
}

type userCar struct { //наша joined структура
	UserId       int64 `db:"userid"` //замапим поля
	Name         string
	Rank         string
	CarId        int64 `db:"carid"`
	Brand        string
	Colour       string
	LicensePlate string
}

func convertJoinedQueryToCar(input userCar) models.Car { //сворачиваем плоскую структуру в нашу рабочую модель
	return models.Car{
		Id:           input.CarId,
		Colour:       input.Colour,
		Brand:        input.Brand,
		LicensePlate: input.LicensePlate,
		Owner: models.User{
			Id:   input.UserId,
			Name: input.Name,
			Rank: input.Rank,
		},
	}
}

func NewCarStorage(pool *pgxpool.Pool) *CarsStorage {
	storage := new(CarsStorage)
	storage.databasePool = pool
	return storage

}

func (storage *CarsStorage) GetCarsList(ctx context.Context, userIdFilter int64, brandFilter string, colourFilter string, licenseFilter string) []models.Car { //TODO test
	query := "SELECT users.id AS userid, users.name, users.rank, c.id AS carid, c.brand, c.colour, c.license_plate FROM users JOIN cars c on users.id = c.user_id WHERE 1=1"

	placeholderNum := 1
	args := make([]interface{}, 0)
	if userIdFilter != 0 { //задаем фильтры через плейсхолдеры
		query += fmt.Sprintf(" AND users.id = $%d", placeholderNum)
		args = append(args, userIdFilter) //сразу же добавляем аргумент для фильтра
		placeholderNum++                  //увеличиваем номер плейсхолдеру
	}
	if brandFilter != "" {
		query += fmt.Sprintf(" AND brand ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", brandFilter))
		placeholderNum++
	}
	if colourFilter != "" {
		query += fmt.Sprintf(" AND colour ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", colourFilter))
		placeholderNum++
	}
	if licenseFilter != "" {
		query += fmt.Sprintf(" AND license_plate ILIKE $%d", placeholderNum)
		args = append(args, fmt.Sprintf("%%%s%%", licenseFilter))
		placeholderNum++
	}

	var dbResult []userCar

	err := pgxscan.Select(ctx, storage.databasePool, &dbResult, query, args...)

	if err != nil {
		log.Errorln(err)
	}

	result := make([]models.Car, len(dbResult))

	for idx, dbEntity := range dbResult {
		result[idx] = convertJoinedQueryToCar(dbEntity) //заполняем результат
	}

	return result
}

func (storage *CarsStorage) GetCarById(ctx context.Context, id int64) models.Car {
	query := "SELECT users.id AS userid, users.name, users.rank, c.id AS carid, c.brand, c.colour, c.license_plate FROM users JOIN cars c on users.id = c.user_id WHERE c.id = $1"

	var result userCar

	err := pgxscan.Get(ctx, storage.databasePool, &result, query, id) //забираем по id

	if err != nil {
		log.Errorln(err)
	}

	return convertJoinedQueryToCar(result)
}

func (storage *CarsStorage) CreateCar(ctx context.Context, car models.Car) error {
	tx, err := storage.databasePool.Begin(ctx) //здесь будем пользоваться транзакцией, чтобы проверка пользователей и вставка нового автомобиля выглядели одним запросом с ее точки зрения
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Errorln(err)
		}
	}()

	query := "SELECT id FROM users WHERE id = $1"

	id := -1

	err = pgxscan.Get(ctx, tx, &id, query, car.Owner.Id)
	if err != nil {
		log.Errorln(err)
		err = tx.Rollback(ctx) //если получили ошибку откатываем транзакцию целиком
		if err != nil {
			log.Errorln(err)
		}
		return err
	}

	if id == -1 {
		return errors.New("user not found")
	}

	insertQuery := "INSERT INTO cars(user_id, colour, brand, license_plate) VALUES ($1,$2,$3,$4)"

	_, err = tx.Exec(ctx, insertQuery, car.Owner.Id, car.Colour, car.Brand, car.LicensePlate) //вызываем exec НЕ У соединения а У транзакции

	if err != nil {
		log.Errorln(err)
		err = tx.Rollback(ctx)
		if err != nil {
			log.Errorln(err)
		}
		return err
	}
	err = tx.Commit(ctx) // в конце посылаем транзакцию, база сохранит значения, если до этого ничего не было откачено
	if err != nil {
		log.Errorln(err)
	}

	return err
}
