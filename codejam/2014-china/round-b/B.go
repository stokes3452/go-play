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
		rectangles := readIntLine()
		Debug("We're dealing with %d rectangles!\n", rectangles)
		var peopleX []int64
		var peopleY []int64
		var peopleWeights []int64
		var realPeopleX []int
		var realPeopleY []int
		var weightedX int64
		var weightedY int64
		for i := 0; i < rectangles; i++ {
			rectangle := getIntsFromLine()
			x1 := rectangle[0]
			y1 := rectangle[1]
			x2 := rectangle[2]
			y2 := rectangle[3]
			peopleX = append(peopleX, x1+x2)
			peopleY = append(peopleY, y1+y2)
			xWidth := (x2 - x1)
			if xWidth < 0 {
				xWidth = xWidth * -1
			}
			yWidth := (y2 - y1)
			if yWidth < 0 {
				yWidth = yWidth * -1
			}
			peopleWeights = append(peopleWeights, (xWidth+1)*(yWidth+1))
			for x := x1; x <= x2; x++ {
				for y := y1; y <= y2; y++ {
					realPeopleX = append(realPeopleX, int(x))
					realPeopleY = append(realPeopleY, int(y))
				}
			}
			weightedX += xWidth * (x1 + x2)
			weightedY += yWidth * (x1 + x2)
		}
		Debug("Our center of weight is %d %d\n", weightedX/int64(len(realPeopleX)), weightedY/int64(len(realPeopleY)))

		totalPeople := len(peopleX)
		var peopleDistances [1000000]int64
		minLocation := 0
		for i := 0; i < totalPeople; i++ {
			for j := i + 1; j < totalPeople; j++ {
				pairX := peopleX[i] - peopleX[j]
				pairY := peopleY[i] - peopleY[j]
				if pairX < 0 {
					pairX = 0 - pairX
				}
				if pairY < 0 {
					pairY = 0 - pairY
				}
				peopleDistances[i] += (pairX + pairY) * peopleWeights[j]
				peopleDistances[j] += (pairX + pairY) * peopleWeights[i]
			}
			if peopleDistances[minLocation] > peopleDistances[i] {
				minLocation = i
			} else if peopleDistances[minLocation] == peopleDistances[i] {
				Debug("We have a dupe!\n")
				if peopleX[minLocation] > peopleX[i] {
					minLocation = i
				} else if peopleY[minLocation] > peopleY[i] {
					minLocation = i
				}
			}
		}

		Debug("We're partying at %d's rectangle!, up in the %d %d house", minLocation, peopleX[minLocation], peopleY[minLocation])

		finalX := peopleX[minLocation] / 2
		finalY := peopleY[minLocation] / 2
		finalDistance := peopleDistances[minLocation] / 2
		/**
		totalPeople = len(realPeopleX)
		finalDistance := 0
		for i := 0; i < totalPeople; i++ {
			pairX := realPeopleX[i] - finalX
			pairY := realPeopleY[i] - finalY
			if pairX < 0 {
				pairX = 0 - pairX
			}
			if pairY < 0 {
				pairY = 0 - pairY
			}
			finalDistance = finalDistance + pairX + pairY
		}*/

		fmt.Fprintf(outputFile, "Case #%d: %d %d %d\n", testNumber, finalX, finalY, finalDistance)
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

func getIntsFromLine() (lineValues []int64) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.ParseInt(string(part), 10, 64)
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
