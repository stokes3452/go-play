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

var fileIn = flag.String("fileIn", "in", "The file to inport")
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

		boxSize := readIntLine()

		if boxSize > 6 {
			panic(fmt.Sprintf("Uh oh, someone's cheating and gave us an over-sized matrix.. length is %d\n", boxSize))
		}
		matrixSize := boxSize * boxSize
		Debug("Matrix size is: %d\n", matrixSize)
		//read inputs.  Build our rows and columns at the same time
		//MAX matrix size is
		var rows [36]sort.IntSlice
		var boxes [36]sort.IntSlice
		valid := true
		for i := 0; i < matrixSize; i++ {
			if valid {
				columns := getIntsFromLine()
				for j := 0; j < matrixSize; j++ {
					Debug("We're going with box %d\n", j/3+(i/3)*3)
					boxes[j/boxSize+(i/boxSize)*boxSize] = append(boxes[j/boxSize+(i/boxSize)*boxSize], columns[j])
					rows[j] = append(rows[j], columns[j])
				}
				valid = valid && validateSlice(columns, matrixSize)
			} else {
				readLine()
			}
		}
		for i := 0; i < matrixSize; i++ {
			if valid {
				valid = valid && validateSlice(rows[i], matrixSize)
				valid = valid && validateSlice(boxes[i], matrixSize)
			}
		}

		if valid {
			fmt.Fprintf(outputFile, "Case #%d: yes\n", testNumber)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: no\n", testNumber)
		}
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func validateSlice(slice sort.IntSlice, matrixSize int) bool {
	sort.Sort(slice)
	for i := 1; i <= matrixSize; i++ {
		if slice[i-1] != i {
			return false
		}
	}
	return true
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
