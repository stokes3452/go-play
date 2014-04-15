package main

import (
	"flag"
	"fmt"
	"bufio"
	"bytes"
	"strconv"
	"os"
	"sort"
)

var fileIn = flag.String("fileIn", "in.txt", "The file to inport")
var fileOut = flag.String("fileOut", "out.txt", "The output file")
var debug = flag.Bool("debug", false, "Debug on or off?")
var reader *bufio.Reader
var outputFile *os.File

func Debug(format string, a ...interface{}) {
	if (*debug) {
		fmt.Printf(format, a...)
	}
}

func init() () {
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
	var rowValues, answer []int
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)

		rowChoice := readIntLine()
		if rowChoice < 1 || rowChoice > 4 {
			panic(fmt.Sprintf("Invalid row choice: %d\n", rowChoice))
		}
		
		Debug("The volunteer chose row: %d\n", rowChoice)
		for currentRow := 1; currentRow <= 4; currentRow++ {
			if currentRow == rowChoice {
				rowValues = getIntsFromLine()
			} else {
				readLine()
			}
		}
		
		rowChoice = readIntLine()
		if rowChoice < 1 || rowChoice > 4 {
			panic(fmt.Sprintf("Invalid row choice: %d\n", rowChoice))
		}
		
		Debug("The volunteer chose row: %d\n", rowChoice)
		for currentRow := 1; currentRow <= 4; currentRow++ {
			if currentRow == rowChoice {
				answer = getAnswerFromLine(rowValues)
			} else {
				readLine()
			}
		}
		
		if len(answer) == 1 {
			fmt.Fprintf(outputFile, "Case #%d: %d\n", testNumber, answer[0])
		} else if len(answer) == 0 {
			fmt.Fprintf(outputFile, "Case #%d: Volunteer cheated!\n", testNumber)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: Bad magician!\n", testNumber)
		}
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func getAnswerFromLine(oldValues []int) (result []int) {
	newValues := getIntsFromLine()
	for _, oldValue := range oldValues {
		for _, newValue := range newValues {
			if oldValue == newValue {
				result = append(result, oldValue)
			}
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