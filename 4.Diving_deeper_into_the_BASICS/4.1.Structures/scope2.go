package main

import "fmt"
import "slurm/client"

func main()  {
	c := client.Client{}

	fmt.Printf("%#v\n", c)
	c.UpdateAvatar()
	fmt.Printf("%#v\n", c)

	c.Name = "Client Name"
	c.Person.Name = "Person name"
	fmt.Println(c.GetName())
	fmt.Println(c.Person.GetName())
}
