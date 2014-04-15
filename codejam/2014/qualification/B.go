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

	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		inputs := getFloatsFromLine()

		factoryCost := inputs[0]
		factoryRate := inputs[1]
		totalTarget := inputs[2]
		currentRate := 2.0
		totalTime := 0.0
		stepTime := 0.0
		for {
			totalTime += stepTime
			timeUntilComplete := totalTarget / currentRate
			timeUntilFactory := factoryCost / currentRate
			stepTime = timeUntilComplete
			Debug("Time until complete: %f\n", timeUntilComplete)
			Debug("Time until factory: %f\n", timeUntilFactory)
			if timeUntilFactory < timeUntilComplete {
				newRate := currentRate + factoryRate
				newTimeUntilComplete := totalTarget / newRate + timeUntilFactory
				if newTimeUntilComplete < timeUntilComplete {
					Debug("Build next factory!\n")
					stepTime = timeUntilFactory
					currentRate = newRate
				} else {
					Debug("We're done, too long!\n")
					break
				}
			} else {
				Debug("We're done, too expensive!\n")
				break
			}
		}
		totalTime += stepTime
		
		fmt.Fprintf(outputFile, "Case #%d: %f\n", testNumber, totalTime)
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