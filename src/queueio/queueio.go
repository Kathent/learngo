package main

import "fmt"

func main() {
	var name1, name string
	fmt.Println("input your name ..")
	fmt.Scanln(&name1, &name)

	fmt.Printf("name1 , name2 %s, %s", name1, name)
}
