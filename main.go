package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	reader := os.Stdin
	wordSlice := convertFileToSlice()
	ui(reader, wordSlice)
}

func convertFileToSlice() [][]byte {
	filepath := "length_5"
	wordSlice := [][]byte{}
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file contents
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var word []byte
	for i := 0; i < len(content); i++ {
		if content[i] == 32 {
			if len(word) > 0 {
				wordSlice = append(wordSlice, word)
				word = nil
			}
		} else {
			word = append(word, content[i])
		}
	}

	if len(word) > 0 {
		wordSlice = append(wordSlice, word)
	}

	return wordSlice
}

func randIndex(len int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len)
	return randomIndex
}

func ui(reader io.Reader, wordSlice [][]byte) {
	scanner := bufio.NewReader(reader)

	word := wordSlice[randIndex(len(wordSlice))]
	for {
		fmt.Println("First word:", string(word))
		fmt.Println("Type 'q' to quit")
		text, err := scanner.ReadString('\n')
		text = strings.TrimSpace(text)
		if err != io.EOF && len(text) > 0 {
			cmd := strings.Split(text, " ")
			fmt.Println(cmd[0])

			if cmd[0] == "q" {
				break
			}

		}
	}
}
