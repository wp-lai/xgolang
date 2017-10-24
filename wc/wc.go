// emulate the unix wc command
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

type stat struct {
	line, word, char int
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		fmt.Println("usage: wc <file1> [<file2> [... <fileN>]]")
		os.Exit(1)
	}

	var lineNumTotal, wordNumTotal, charNumTotal int
	ch := make(chan stat)

	for _, filename := range args[1:] {
		go func(filename string) {
			lineNum, wordNum, charNum, err := count(filename)
			if err != nil {
				fmt.Printf("error count: %s", err)
			}

			fmt.Printf("\t%d\t%d\t%d\t%s\n", lineNum, wordNum, charNum, filename)
			ch <- stat{lineNum, wordNum, charNum}
		}(filename)
	}

	for range args[1:] {
		num := <-ch
		lineNumTotal += num.line
		wordNumTotal += num.word
		charNumTotal += num.char
	}
	fmt.Printf("\t%d\t%d\t%d\t%s\n", lineNumTotal, wordNumTotal, charNumTotal, "total")

}

// count returns the number of lines, words, and chars in the file
func count(filename string) (int, int, int, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return 0, 0, 0, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return 0, 0, 0, err
	}
	defer f.Close()

	r := regexp.MustCompile("[^\\s]+")
	var wordNum, lineNum, charNum int

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		lineNum++
		for range r.FindAllString(line, -1) {
			wordNum++
		}
		charNum += len(line)
	}

	return lineNum, wordNum, charNum, nil
}
