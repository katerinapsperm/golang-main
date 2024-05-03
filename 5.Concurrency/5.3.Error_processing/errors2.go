package main

import (
	"errors"
	"fmt"
	"slurm/client"
)

type ServerError struct {
	err error
}

func (s *ServerError) Error() string  {
	return s.err.Error()
}

func NewServerError(msg string) error {
	return &ServerError{err: errors.New(msg)}
}

var internalErr = NewServerError("внутренняя ошибка")

func getClient() (client.Client, error) {
	cl := client.Client{}
	err := fmt.Errorf("ошибка получения клиента: %v", internalErr)
	return cl, err
}

func main()  {
	res, err := getClient()
	if err != nil {
		if errors.Is(err, internalErr) {
			fmt.Println("произошла внутренняя ошибка")
		}
		newErr := errors.Unwrap(err)
		fmt.Println(newErr)
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)
}