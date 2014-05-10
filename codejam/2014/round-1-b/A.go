package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var fileIn = flag.String("fileIn", "in.txt", "The file to inport")
var fileOut = flag.String("fileOut", "out.txt", "The output file")
var debug = flag.Bool("debug", false, "Debug on or off?")
var reader *bufio.Reader
var outputFile *os.File

func Debug(format string, a ...interface{}) {
	if *debug {
		fmt.Printf(format, a...)
	}
}

func init() {
	flag.Parse()
	Debug("Debug is %t\n", *debug)
	Debug("File in is: %s\n", *fileIn)
	Debug("File out is: %s\n", *fileOut)

	file, err := os.Open(*fileIn)
	if err != nil {
		panic(fmt.Sprintf("File read failed with err: %s", err))
	}

	reader = bufio.NewReaderSize(file, 16000)
	outputFile, err = os.Create(*fileOut)
	if err != nil {
		panic(fmt.Sprintf("File creation failed with err: %s", err))
	}
}

func main() {
	testsRemaining := readIntLine()
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		numStrings := readIntLine()
		var gameStrings [][]byte
		var currentCondensed []byte
		winning := true
		for i := 0; i < numStrings; i++ {
			currentString := readLine()
			condensedString := getCondensedString(currentString)
			if len(currentCondensed) > 0 {
				if !bytes.Equal(condensedString, currentCondensed) {
					winning = false
				}
			} else {
				currentCondensed = condensedString
			}
			//force a memory copy so we don't leak over our buffers..
			newSlice := make([]byte, len(currentString))
			copy(newSlice, currentString)
			gameStrings = append(gameStrings, newSlice)
			Debug("Current string: %s which condenses to %s\n", currentString, condensedString)
		}

		if winning {
			moves := calculateMoves(gameStrings, currentCondensed)
			fmt.Fprintf(outputFile, "Case #%d: %d\n", testNumber, moves)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: Fegla Won\n", testNumber)
		}

		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func calculateMoves(gameStrings [][]byte, condensedString []byte) (moves int) {
	//for each character in our final target string
	for _, condensedChar := range condensedString {
		var stringDeltas sort.IntSlice
		//grab a count of how many times that character appears in our input strings
		//in the end, we want to go for the MEDIAN character count
		for gameIndex, currentString := range gameStrings {
			var index int
			var currentChar byte
			for index, currentChar = range currentString {
				if currentChar != condensedChar {
					gameStrings[gameIndex] = currentString[index:]
					index = index - 1
					break
				}
			}
			Debug("We're cutting off at index %d\n", index)
			stringDeltas = append(stringDeltas, index)
		}
		sort.Sort(stringDeltas)
		guessInt := stringDeltas[len(gameStrings)/2]
		for _, stringDelta := range stringDeltas {
			if stringDelta > guessInt {
				moves += stringDelta - guessInt
			} else {
				moves += guessInt - stringDelta
			}
		}
	}

	return
}

func getCondensedString(fullString []byte) (result []byte) {
	var currentChar byte
	for _, char := range fullString {
		if char != currentChar {
			currentChar = char
			result = append(result, char)
		}
	}
	return
}

func readLine() (line []byte) {
	line, prefix, err := reader.ReadLine()
	if err != nil {
		panic("Error reading line for test")
	}
	if prefix {
		panic("Did not complete read.  Increase buffer size and try try again")
	}
	return
}

func getFloatsFromLine() (lineValues sort.Float64Slice) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.ParseFloat(string(part), 64)
		if err != nil {
			panic(fmt.Sprintf("bad float provided: %s", part))
		}
		lineValues = append(lineValues, currentValue)
	}
	return
}

func getIntsFromLine() (lineValues sort.IntSlice) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.Atoi(string(part))
		if err != nil {
			panic(fmt.Sprintf("bad int provided: %s", part))
		}
		lineValues = append(lineValues, currentValue)
	}
	return
}

func readIntLine() (response int) {
	line := readLine()
	response, err := strconv.Atoi(string(line))
	if err != nil {
		panic("Err reading int value from line")
	}
	return
}
