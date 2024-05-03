package main

import "fmt"

//type Avatar struct {
//	URL string
//	Size int64
//}
//
//type Client struct {
//	ID int64
//	Name string
//	Age int
//	IMG *Avatar
//}

func main()  {
	i := new(int)
	_ = i

	client := Client{}
	//client.IMG = new(Avatar)
	client.IMG = Avatar{}
	//fmt.Printf("%#v\n", client)
	// main.Client{ID:0, Name:"", Age:0, IMG:main.Avatar{URL:"", Size:0}}
	// main.Client{ID:0, Name:"", Age:0, IMG:(*main.Avatar)(nil)}

	//client.IMG = new(Avatar)
	updateAvatar(&client)
	fmt.Printf("%#v\n", client)
	//fmt.Printf("%#v\n", client.IMG)

	updateClient(client)
	//fmt.Printf("%#v\n", client)
}

func updateAvatar(client *Client) {
	client.IMG.URL = "updated_url"
}

func updateClient(client Client) {
	client.Name = "Артем"
}
