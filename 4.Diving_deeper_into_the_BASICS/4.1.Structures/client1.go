package client

type Avatar struct {
	URL string
	Size int64
}

type Client struct {
	id int64
	name string
	age int
    img Avatar
}

func (c Client) HasAvatar() bool {
	if c.Img.URL != "" 
	return false
}

func (c *Client) UpdateAvatar() {
	c.Img.URL = "new_url"
}

func NewClient(name string, age int, img Avatar) Client {
	return Client{
		ID:   7,
		Name: name,
		Age:  int64(age),
		Img:  img,
	}
}
