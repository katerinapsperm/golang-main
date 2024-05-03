package main

import "fmt"

type Caller interface {
	Call(number int) error
}

type Sender interface {
	Send(msg string) error
}

type IPhone interface {
	Caller
	Sender
	MyFunc()
	MyFunc2()
}

type Email struct {
	Address string
}

func (e *Email) Send(msg string) error {
	fmt.Printf("Сообщение \"%v\" отправлено по почете на адрес %v \n", msg, e.Address)
	return nil
}

type Phone struct {
	Number int
	Balance int
}

func (p *Phone) Send(msg string) error {
	fmt.Printf("Сообщение \"%v\" отправлено по телефону с номера %v \n", msg, p.Number)
	return nil
}

func Notify(i interface{}) {
	switch i.(type) {
	case int:
		fmt.Println("Число не поддерживается")
	}

	s, ok := i.(Sender)
	if !ok {
		fmt.Println("Ошибка утверждения интерфейса")
		return
	}

	err := s.Send("Сообщение пустого интерфейса")
	if err != nil {
		fmt.Println("Ошибка")
		return
	}
	fmt.Println("Success")
}

func main()  {
	email := &Email{"dev@slurm.io"}
	Notify(email)

    phone := &Phone{Number: 7777, Balance: 10}
	Notify(phone)

    Notify(2)
	Notify("Строка")

	sl := [5]int64{1,2,3,4,5}
	Notify(sl)
}
