package main

import "fmt"

func main() {
	s := make([]int, 3, 6)

	s2 := s[1:2]

	fmt.Println(len(s2))
	fmt.Println(cap(s2))
}
