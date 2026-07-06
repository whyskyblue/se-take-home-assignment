package main

import (
	"bufio"
	"fmt"
	"mcdonalds-order-controller/service"
	"os"
	"strings"
	"time"
)

func main() {
	controller, err := service.NewOrderController("scripts/result.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer controller.Close()

	service.PrintHelp()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch input {
		case "n":
			controller.AddNormalOrder()
		case "v":
			controller.AddVIPOrder()
		case "+":
			controller.AddBot()
		case "-":
			controller.RemoveBot()
		case "s":
			controller.LogStatus()
		case "h":
			service.PrintHelp()
		case "q", "quit", "exit":
			controller.Log("Shutting down...")
			time.Sleep(25 * time.Second)
			return
		default:
			fmt.Println("Unknown command. Type 'h' for help.")
		}
	}
}
