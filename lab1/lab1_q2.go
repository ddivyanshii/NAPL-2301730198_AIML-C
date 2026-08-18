package main

import (
	"fmt"
)

func clearInput() {
	var trash string
	fmt.Scanln(&trash)
}

func main() {
	for {
		var choice int

		fmt.Println("\nSimple Calculator")
		fmt.Println("1. Integer")
		fmt.Println("2. Float")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			clearInput()
			fmt.Println("Invalid choice, returning to menu.")
			continue
		}

		if choice == 3 {
			break
		}

		if choice == 1 {
			var a, b int
			fmt.Print("Enter 1st integer: ")
			_, errA := fmt.Scan(&a)
			if errA != nil {
				clearInput()
				fmt.Println("Invalid input, returning to menu.")
				continue
			}

			fmt.Print("Enter 2nd integer: ")
			_, errB := fmt.Scan(&b)
			if errB != nil {
				clearInput()
				fmt.Println("Invalid input, returning to menu.")
				continue
			}

			fmt.Println("Integer Addition:", a+b)
			fmt.Println("Integer Subtraction:", a-b)
			fmt.Println("Integer Multiplication:", a*b)
		}

		if choice == 2 {
			var x, y float64
			fmt.Print("Enter 1st decimal number: ")
			_, errX := fmt.Scan(&x)
			if errX != nil {
				clearInput()
				fmt.Println("Invalid input, returning to menu.")
				continue
			}

			fmt.Print("Enter 2nd decimal number: ")
			_, errY := fmt.Scan(&y)
			if errY != nil {
				clearInput()
				fmt.Println("Invalid input, returning to menu.")
				continue
			}

			fmt.Println("Float Addition:", x+y)
			fmt.Println("Float Subtraction:", x-y)
			fmt.Println("Float Multiplication:", x*y)
		}
	}
}
