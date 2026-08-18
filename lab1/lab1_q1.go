package main

import "fmt"

func main() {

	var a, b int
	fmt.Print("Enter 1st integer: ")
	fmt.Scan(&a)
	fmt.Print("Enter 2nd integer: ")
	fmt.Scan(&b)

	fmt.Println("Integer Addition:", a+b)
	fmt.Println("Integer Subtraction:", a-b)
	fmt.Println("Integer Multiplication:", a*b)

	var x, y float64
	fmt.Printf("Enter 1st decimal number: ")
	fmt.Scan(&x)
	fmt.Printf("Enter 2nd decimal number: ")
	fmt.Scan(&y)

	fmt.Println("Float Addition:", x+y)
	fmt.Println("Float Subtraction:", x-y)
	fmt.Println("Float Multiplication:", x*y)

}
