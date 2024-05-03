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

func Notify(s Sender{}) {
	err := s.Send("Notify message")
	if err != nil {
		fmt.Println("err")
		return
	}

    swith s.(type) {
    case *Email:
        fmt.Println("Сообщение по почте")
    case *Phone:
        phone := s.(*Phone)
        fmt.Println(phone.Balance)
    }
    
	fmt.Println("Success")
}

func main()  {
	email := &Email{"dev@slurm.io"}
	Notify(email)

    phone := &Phone{Number: 7777, Balance: 10}
	Notify(phone)
}
