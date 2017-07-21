package main

import (
	"os"
	"bufio"
	"fmt"
	"io"
)

func main() {
	//readFile()
	writeFile()
}
func writeFile() {
	file, err := os.OpenFile("aa.txt", os.O_APPEND, 066)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("\r\n new line...")
	writer.WriteString("\r\n new line again..")
	writer.Flush()
}
func readFile() {
	open, error := os.Open("aa.txt")
	if error != nil {
		panic(error)
	}

	defer open.Close()

	reader := bufio.NewReader(open)
	for {
		string, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}

		fmt.Println(string)

		if err == io.EOF {
			return
		}
	}

	open.Close()
}
