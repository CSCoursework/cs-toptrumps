package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var scanner = bufio.NewScanner(os.Stdin)

// Text prompts the user for a line of text, which it then returns.
func Text(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

// Int prompts the user for an input, which is validates to ensure it is an integer, and then returns as an integer
func Int(prompt string) int {
	for {
		n, err := strconv.ParseInt(Text(prompt), 10, 8) // Parse the string as a base-10, 8-bit number
		if err != nil {
			fmt.Println("Number is not an integer.")
		} else {
			return int(n) // strconv.ParseInt returns an int64
		}
	}
}

// Options presents a collection of options to the user, and invites them to choose one. The index of that option and the
// option value are then returned.
func Options(prompt string, items []string) (int, string) {
	for {
		fmt.Println(prompt)
		for i, v := range items {
			fmt.Printf("  %d: %s\n", i+1, v)
		}
		num := Int("> ") - 1
		if !(num >= 0 && num < len(items)) {
			fmt.Println("Out of bounds")
		} else {
			return num, items[num]
		}
	}
}
