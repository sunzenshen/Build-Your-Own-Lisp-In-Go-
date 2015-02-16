package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sunzenshen/golispy/mpcinterface"
)

func main() {
	// Version and Exit Information
	fmt.Println("Lispy Version 0.0.0.0.3")
	fmt.Print("Press Ctrl+c to Exit\n\n")
	// For reading lines of user input
	scanner := bufio.NewScanner(os.Stdin)

	lispy := mpcinterface.InitLispy()
	defer mpcinterface.CleanLispy(lispy)

	for {
		// Prompt
		fmt.Print("lispy> ")
		// Read a line of user input
		scanner.Scan()
		input := scanner.Text()
		// Echo input back to user
		lispy.ReadEvalPrint(input)
	}
}
