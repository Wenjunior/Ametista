package print

import (
	"os"
	"fmt"
	"bufio"
)

import (
	"amt/utils/colors"
)

func eprintln(err string, color string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", color, err, colors.RESET)
}

func Eprintln(err string) {
	eprintln(err, colors.YELLOW)
}

func Panic(err error) {
	eprintln(err.Error(), colors.RED)

	os.Exit(1)
}

func BufferedPrint(items []string) {
	writer := bufio.NewWriter(os.Stdout)

	writer.WriteString(colors.GREEN)

	for _, item := range items {
		writer.WriteString(item + "\n")
	}

	writer.WriteString(colors.RESET)

	writer.Flush()
}

func Cprintln(msg string, color string) {
	fmt.Printf("%s%s%s\n", color, msg, colors.RESET)
}