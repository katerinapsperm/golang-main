package client

type Person struct {
	Name string
	Age int
}

func (p Person) GetName() string {
	return p.Name
}

type Avatar struct {
	URL string
	Size int64
}

type Client struct {
	ID int64
	Img Avatar
	Name string
	Age int64
}

func (c Client) HasAvatar() bool {
	if c.Img.URL != "" {
		return true
	}
	return false
}

func (c *Client) UpdateAvatar() {
	c.Img.URL = "new_url"
}

func (c Client) GetName() string {
	return c.Name
}

func NewClient(name string, age int, img Avatar) *Client {
	return &Client{
		ID:   7,
		Name: name,
		Age:  int64(age),
		Img:  img,
	}
}
