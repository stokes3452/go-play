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
		interestingStuff := getIntsFromLine()
		var zipCodes sort.IntSlice
		for i := 0; i < interestingStuff[0]; i++ {
			zipCodes = append(zipCodes, readIntLine())
		}

		flights := make(map[int]map[int]int)
		for i := 0; i < interestingStuff[1]; i++ {
			//source, destination
			flight := getIntsFromLine()
			source := zipCodes[flight[0]-1]
			destination := zipCodes[flight[1]-1]
			if flights[source] == nil {
				flights[source] = make(map[int]int)
			}
			flights[source][destination] = destination
			if flights[destination] == nil {
				flights[destination] = make(map[int]int)
			}
			flights[destination][source] = source
		}

		Debug("Read in %d flight options!\n", len(flights))
		sort.Sort(zipCodes)

		var flightPlan []int
		var index int
		var zip int
		for index, zip = range zipCodes {
			Debug("Trying out zip %d\n", zip)
			flightPlan = getFlightPlan(removeSliceElement(zipCodes, index), []int{zip}, flights)
			if flightPlan != nil {
				break
			}
		}

		fmt.Fprintf(outputFile, "Case #%d: %d", testNumber, zip)
		for index = len(flightPlan) - 1; index >= 0; index-- {
			fmt.Fprintf(outputFile, "%d", flightPlan[index])
		}
		fmt.Fprintf(outputFile, "\n")

		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func getMinimumSpanningTreeIndex(remainingZipCodes, currentZipCodes []int, flights map[int]map[int]int) (minimumIndex int) {
	remainingZipMap := make(map[int]bool)
	var zip int
	//re-organize into a map for easy existance-checking
	for _, zip = range remainingZipCodes {
		remainingZipMap[zip] = true
	}
	//no need to initialize anything into our channel - the default serves 2 purposes here
	ch := make(chan int, len(remainingZipCodes))
	//as long as we have anything left to find a path to, search on!
	for len(remainingZipMap) > 0 {
		select {
		case zip = <-ch:
			for _, newZip := range flights[zip] {
				_, ok := remainingZipMap[newZip]
				if ok {
					ch <- newZip
					delete(remainingZipMap, newZip)
				}
			}
			break
		default:
			ch <- currentZipCodes[minimumIndex]
			minimumIndex++
			break
		}
	}
	return minimumIndex - 1
}

func getFlightPlan(remainingZipCodes, currentZipCodes []int, flights map[int]map[int]int) (flightPlan []int) {
	//If we're done, return
	if len(remainingZipCodes) == 0 {
		return
	}
	//find the minimum index that we can use, without blocking ourselves from building a spanning tree
	minimumZipIndex := getMinimumSpanningTreeIndex(remainingZipCodes, currentZipCodes, flights)
	Debug("Min index: %d\n", minimumZipIndex)

	//for each sorted zip code
	//IF we can get to it, search on!
	for remainingIndex, targetZip := range remainingZipCodes {
		for oldIndex := len(currentZipCodes) - 1; oldIndex >= minimumZipIndex; oldIndex-- {
			oldZip := currentZipCodes[oldIndex]
			//IF WE FIND OUR TARGET
			_, ok := flights[oldZip][targetZip]
			if ok {
				goto continuation
			}
			continue
		continuation:
			newRemainingZipCodes := removeSliceElement(remainingZipCodes, remainingIndex)
			//currentZipCodes = start..index
			newZipCodes := make([]int, oldIndex+1)
			copy(newZipCodes, currentZipCodes[0:oldIndex+1])
			if len(newRemainingZipCodes) > 0 {
				newPlan := getFlightPlan(newRemainingZipCodes, append(newZipCodes, targetZip), flights)

				if newPlan != nil {
					Debug("Got a valid one! %d\n", targetZip)
					return append(newPlan, targetZip)
				} else {
					Debug("Backing up a level!")
				}
			} else {
				return []int{targetZip}
			}

		}
	}
	return nil
}

func removeSliceElement(oldSlice []int, index int) (newSlice []int) {
	newSlice = make([]int, len(oldSlice)-1)
	if index > 0 {
		copy(newSlice[0:index], oldSlice[0:index])
	}
	if index < len(oldSlice)-1 {
		copy(newSlice[index:], oldSlice[index+1:])
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
