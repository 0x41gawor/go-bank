package main

import "fmt"

type Arabia func(a, b string)

func makeAsiaFromArabia(f apiArabia) Arabia {
	return func(a, b string) {
		fmt.Printf("Elo----: %s\n", a+b)
		if err := f(a, b); err != nil {

		}
		fmt.Printf("Elo-----: %s\n", a+b)
	}
}

type apiArabia func(a, b string) error

func apiA(a, b string) error {
	fmt.Printf("Jol: %s\n", a+b)
	return nil
}

func main() {
	// server := NewApiServer(":3000")
	// server.Run()
	// fmt.Println("Yeah buddy!!")
	nowa := makeAsiaFromArabia(apiA)
	nowa("xx", "zz")
}
