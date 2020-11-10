package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var scanner = bufio.NewScanner(os.Stdin)

func Text(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

func promptInt(prompt string, bitsize int) int64 {
	for {
		n, err := strconv.ParseInt(Text(prompt), 10, bitsize)
		if err != nil {
			fmt.Println("Number is not an integer.")
		} else {
			return n
		}
	}
}

func Int(prompt string) int {
	return int(promptInt(prompt, 8))
}

func Options(prompt string, items []string) (int, string) {
	for {
		fmt.Println(prompt)
		for i, v := range items {
			fmt.Printf("  %d: %s\n", i+1, v)
		}
		num := Int("> ") - 1
		if ! (num >= 0 && num < len(items)) {
			fmt.Println("Out of bounds")
		} else {
			return num, items[num]
		}
	}
}
