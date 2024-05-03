package main

import "fmt"
import "slurm/client"

//type Avatar struct {
//	URL string
//	Size int64
//}
//
//type Client struct {
//	ID int64
//	Name string
//	Age int
//	IMG Avatar
//}
//
//func (c Client) HasAvatar() bool {
//	if c.IMG.URL != "" {
//		return true
//	}
//	return false
//}
//
//func (c *Client) UpdateAvatar() {
//	c.IMG.URL = "new_url"
//}

func main()  {
	client := Client{
		ID:   7,
		Name: "Андрей",
		Age:  20,
		IMG: Avatar{
			URL: "some_url",
			Size: 10,
		},
	}
	fmt.Printf("%#v\n", client)
	client.UpdateAvatar()
	fmt.Printf("%#v\n", client)
}
