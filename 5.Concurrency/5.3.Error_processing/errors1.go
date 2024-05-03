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

func getClient() (client.Client, error) {
	cl := client.Client{}
	err := fmt.Errorf("ошибка получения клиента: %w", internalErr)
	return cl, err
}

func main()  {
	res, err := getClient()
	if err != nil {
		fmt.Println(err.Error())
        return
		}
	fmt.Println(res)
}
