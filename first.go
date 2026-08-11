package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!") //print hello go
	fmt.Println("Divyanshi,2301730198")

	var a, b float64

	var choice int
	// Validate numeric input
	for {
		fmt.Print("Enter two numbers (0-100): ")
		if _, err := fmt.Scan(&a, &b); err != nil {
			fmt.Println("Invalid input! Enter numbers only.")
			fmt.Scanln()
			continue
		}
		if a < 0 || a > 100 || b < 0 || b > 100 {
			fmt.Println("Values must be between 0 and 100.")
			continue
		}
		break
	}
	// Validate operation
	for {
		fmt.Print("1.Add  2.Subtract  3.Multiply: ")
		fmt.Scan(&choice)
		if choice >= 1 && choice <= 3 {
			break
		}
		fmt.Println("Invalid choice!")
	}
	if choice == 1 {
		fmt.Println("Result:", a+b)
	} else if choice == 2 {
		fmt.Println("Result:", a-b)
	} else {
		fmt.Println("Result:", a*b)
	}

	fmt.Print("Enter two integers: ")
	fmt.Scan(&a, &b)
	fmt.Println("Integer Addition:", a+b)
	fmt.Println("Integer Subtraction:", a-b)
	fmt.Println("Integer Multiplication:", a*b)
	// Floating-point operations
	var x, y float64
	fmt.Print("Enter two decimal numbers: ")
	fmt.Scan(&x, &y)
	fmt.Println("Float Addition:", x+y)
	fmt.Println("Float Subtraction:", x-y)
	fmt.Println("Float Multiplication:", x*y)
}