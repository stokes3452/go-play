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

type pair map[string] bool

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
		readLine()
		bookValues := getIntsFromLine()
		oddValues := sort.IntSlice{}
		evenValues := sort.IntSlice{}
		valueMap := []bool{}
		for _, bookValue := range bookValues {
			compareValue := bookValue
			if bookValue < 0 {
				compareValue = bookValue * -1
			}
			if compareValue % 2 == 1 {
				oddValues = append(oddValues, bookValue)
				valueMap = append(valueMap, true)
			} else {
				evenValues = append(evenValues, bookValue)
				valueMap = append(valueMap, false)
			}
		}
		sort.Sort(oddValues)
		sort.Sort(evenValues)
		oddIndex := 0
		evenIndex := len(evenValues) - 1
		fmt.Fprintf(outputFile, "Case #%d:", testNumber)
		for _, value := range valueMap {
			if value {
				//odd!
				fmt.Fprintf(outputFile, " %d", oddValues[oddIndex])
				oddIndex++
			} else {
				fmt.Fprintf(outputFile, " %d", evenValues[evenIndex])
				evenIndex--
			}
		}
		
		fmt.Fprintf(outputFile, "\n")
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
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
		panic(fmt.Sprintf("Err reading int value from line: %d", line))
	}
	return
}