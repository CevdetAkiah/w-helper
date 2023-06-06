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

	ui(reader, wordSlice, input())
}

func ui(reader io.Reader, wordSlice [][]byte, input chan string) {
	quit := false
	// word := wordSlice[randIndex(len(wordSlice))]
	word := []byte{98, 108, 111, 110, 100}
	index := 0
	letterMap := make(map[int]byte)
	tell("Type 'q' to quit")
	tell("Type 'n' for another word")
	tell("Type 'y' for yes")
	tell("First word:", string(word))
	tell("Did you score any green hits?")

	for !quit {
		userInput := <-input
		cmd := strings.Split(userInput, " ")

		switch cmd[0] {
		case "q":
			quit = true
		case "n":
			word = wordSlice[randIndex(len(wordSlice))]
			tell(string(word))
			tell("Did you score any green hits?")
			continue
		case "y":
			tell("new commands: v for yes, x for no")
			testWord(string(word), wordSlice, index)
		case "v":
			// record the letter and index
			letterMap[index] = word[index]
			tell("letters stored so far:")
			for i, v := range letterMap {
				fmt.Printf("Letter: %s  -  stored at index: %d\n", string(v), i+1)
			}
			// search for the next hit
			tell("Type p to search for the next letter")
			wordSlice = newWordSlice(index, wordSlice, letterMap)

			for _, v := range wordSlice {
				fmt.Println(string(v))
			}

		case "x":
			index++
			testWord(string(word), wordSlice, index)
		}
	}
}

func newWordSlice(index int, wordSlice [][]byte, letterMap map[int]byte) [][]byte {
	fmt.Println(string(letterMap[index]))
	validCount := 0

	for _, v := range wordSlice {
		if len(v) > index && v[index] == letterMap[index] {
			wordSlice[validCount] = v
			validCount++
		}
	}

	wordSlice = wordSlice[:validCount]
	return wordSlice
}

func testWord(word string, wordSlice [][]byte, index int) {
	test := ""
	i := index
	fmt.Println(i)
	if i >= 5 {
		return
	}
	for _, v := range wordSlice {
		if word[i] == v[i] {
			match := true
			for j := 0; j < 5; j++ {
				if j == i {
					j++
				} else if word[j] == v[j] {
					match = false
					break
				}
			}
			if match {
				test = string(v)
				break
			}
		}
	}
	tell("Any green hits? ", string(test))
}

func tell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println(toGUI)
}

func input() chan string {
	line := make(chan string)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if err != io.EOF && len(text) > 0 {
				line <- text
			}
		}
	}()
	return line
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
