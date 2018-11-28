package main

import "fmt"
import "os"

type Staff struct {
	id   int
	name string
}

func (s Staff) Print() {
	fmt.Println(s.id)
	fmt.Println(s.name)
}

func staff_new(id int, name string) Staff {
	s := Staff {}

	s.id = id;
	s.name = name;

	return s
}

func main() {
	employees := make(map[int]Staff)
	staff := staff_new(123, "Al")

	employees[123] = staff

	staff.Print()

	for _, s := range employees {
                s.Print()
	}
	
	os.Exit(0)
}
