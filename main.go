package main

import (
	"encoding/json"
	"fmt"
	"github.com/GianlucaGuarini/OSXKeyboard"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const MAX_LETTER_COUNT = 64

// load the blacklist from the json file
func loadBlacklist(path string) []string {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	// defer the close of the file reader
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []string
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

// true if the letter pressed is a command key
func isCommandKey(letter string) bool {
	return len(letter) > 1
}

// async callback called on any key press
func onKeyPress() func(string) {
	var chars []string

	blacklist := loadBlacklist(blacklistFilePath())

	// remove the first char from the chars list
	shift := func() {
		if len(chars) > 0 {
			chars = chars[1:]
		}
	}

	// clean up the chars
	clear := func() {
		chars = chars[:0]
	}

	return func(letter string) {
		// command keys like ENTER, ESC... will be skipped
		// notice that the SPACE will be skipped as well
		if isCommandKey(letter) {
			if letter == "DELETE" {
				shift()
			}

			return
		}

		// append the new letter to the queue
		chars = append(chars, strings.ToLower(letter))

		// remove the first letter if the length is longer than the max amount
		if len(chars) > MAX_LETTER_COUNT {
			shift()
		}

		// sleep because i got ya
		if isInBlacklist(strings.Join(chars, ""), blacklist) {
			clear()
			sleep()
		}
	}
}

// check if the letters typed match one of the words in the blacklist
func isInBlacklist(queue string, blacklist []string) bool {
	for _, word := range blacklist {
		if strings.Contains(queue, word) {
			return true
		}
	}

	return false
}

// let me sleep...
func sleep() {
	cmd := exec.Command("pmset", "sleepnow")
	cmd.Run()
}

// generate the path to the blacklist file
func blacklistFilePath() string {
	if len(os.Args) < 2 {
		panic("Please provide the path to your blacklist json file")
	}

	path := os.Args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	return path
}

func main() {
	go OSXKeyboard.Listen()

	OSXKeyboard.Subscribe(onKeyPress())

	select {} // block forever
}
