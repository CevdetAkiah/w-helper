package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	word := []byte{122, 105, 110, 99, 121}
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
			tell("new commands: v for yes, x for no, z to test a new word in the same index followed by how many times this has failed eg 'z 2'")
			testWord(string(word), wordSlice, index, []string{}, letterMap)
		case "v":
			// record the letter and index
			letterMap[index] = word[index]
			tell("letters stored so far:")
			for i, v := range letterMap {
				fmt.Printf("Letter: %s  -  stored at index: %d\n", string(v), i+1)
			}
			// search for the next hit
			wordSlice = newWordSlice(index, wordSlice, letterMap)
			word = wordSlice[randIndex(len(wordSlice))]
			tell("new word: ", string(word))
			index++

			// for _, v := range wordSlice {
			// 	fmt.Println(string(v))
			// }

		case "x":
			index++

			testWord(string(word), wordSlice, index, []string{}, letterMap)
		case "z":
			testWord(string(word), wordSlice, index, cmd, letterMap)

		}

	}
}

func newWordSlice(index int, wordSlice [][]byte, letterMap map[int]byte) [][]byte {
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

func testWord(word string, wordSlice [][]byte, index int, cmd []string, letterMap map[int]byte) {
	test := ""
	i := index

	foundLetters := len(letterMap)
	var skipCount int
	var err error
	if len(cmd) != 0 && len(cmd) >= 2 {
		skipCount, err = strconv.Atoi(cmd[1]) // cyle past the words we have already used
		if err != nil {
			fmt.Println("Please type fail count as an integer, error: ", err)
		}
	}

	if i >= 5 {
		return
	}

	for _, v := range wordSlice {
		if word[i] == v[i] {
			match := true
			if skipCount > 0 {
				skipCount--
				match = false
			}
			for j := 0; j < 5; j++ {
				if j == i {
					j++
				} else if word[j] == v[j] {
					if foundLetters > 0 {
						foundLetters--
						continue
					} else {
						match = false
						break
					}
				}
			}
			if match {
				test = string(v)
				break
			}
		}
	}
	// if len(string(test)) == 0 {
	// 	for _, v := range wordSlice {
	// 		if v[i] == word[i] {
	// 			fmt.Println(string(v))
	// 		}
	// 	}
	// }
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
