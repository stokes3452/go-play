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

	var lawn [100][100]int
	var biggestRow [100]int
	var biggestColumn [100]int
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		dimensions := getIntsFromLine()
		for x := 0; x < dimensions[1]; x++ {
			biggestColumn[x] = 0
		}
		//go ahead and pull out the biggest value from each row while we're here
		for y := 0; y < dimensions[0]; y++ {
			heights := getIntsFromLine()
			//100 is our max_height, so 1 up it
			biggestRow[y] = 0
			for x := 0; x < dimensions[1]; x++ {
				lawn[y][x] = heights[x]
				if heights[x] > biggestRow[y] {
					biggestRow[y] = heights[x]
				}
				if heights[x] > biggestColumn[x] {
					biggestColumn[x] = heights[x]
				}
			}
		}

		valid := true
		//for each X, Y
		//If X is not the biggest in its row OR its column, fail
		for x := 0; x < dimensions[1]; x++ {
			//First, find the things items in each row that are not the biggest
			for y := 0; y < dimensions[0]; y++ {
				lawnHeight := lawn[y][x]
				if lawnHeight < biggestRow[y] && lawnHeight < biggestColumn[x] {
					valid = false
					Debug("Invalid because of %d,%d\n", y, x)
				}
			}
		}
		if valid {
			fmt.Fprintf(outputFile, "Case #%d: YES\n", testNumber)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: NO\n", testNumber)
		}
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
		panic("Err reading int value from line")
	}
	return
}
