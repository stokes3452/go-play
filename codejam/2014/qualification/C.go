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
		
		values := getIntsFromLine()
		if len(values) != 3 {
			panic(fmt.Sprintf("Invalid number of input values: %d", len(values)))
		}
		
		rows := values[0]
		columns := values[1]
		mines := values[2]

		gameSize := rows * columns
		emptyTiles := gameSize - mines
		Debug("%d rows, %d columns, and %d mines\n", rows, columns, mines)
		
		fmt.Fprintf(outputFile, "Case #%d:\n", testNumber)
		if mines >= gameSize {
			fmt.Fprintf(outputFile, "Impossible\n")
			Debug("Trivially impossible, too many mines\n")
		} else if rows == 1 || columns == 1 {
			getComplexSolution(outputFile, emptyTiles, rows, columns)
		} else {
			if emptyTiles == 1 {
				Debug("Simple solution time!")
				getSimpleSolution(outputFile, mines, rows, columns)
			} else if mines == 0 || emptyTiles >= 4 {
				getComplexSolution(outputFile, emptyTiles, rows, columns)
			} else {
				fmt.Fprintf(outputFile, "Impossible\n")
				Debug("No one expects the 2- or 3- click solution!\n")
			}
		}
	}
}

func getComplexSolution(outputFile *os.File, emptyTiles, rows, columns int) {
	//. empty
	//c clicked
	//* has a mine
	
	//In the normal (2x2 or greater) use-case, we don't want to ever exceed a maximum row length of 2
	fullRowLength := emptyTiles / 2
	//However, if we only have 1 row or 1 column, then we can use the full row length
	if rows == 1 || columns == 1 {
		fullRowLength = emptyTiles
	}
	
	//We also clearly can't print more things per row than we have columns
	if fullRowLength > columns {
		fullRowLength = columns
	}
	
	//Calculate the number of full rows that we can stamp out, along with the remainder
	fullRows := emptyTiles / fullRowLength
	remainingCells := emptyTiles % fullRowLength
	Debug("For %dx%d, we have a starting %d full rows of %d length, with a remainder of %d\r\n", rows, columns, fullRows, fullRowLength, remainingCells)
	
	//If we have 2 or more remaining cells we're fine.  If we have 1, we have to make sure we can overflow bonus remainder from the other full rows..
	partialRows := 0
	if remainingCells == 1 {
		//If we have enough rows to only steal 1 mine, we're good (have to maintain > 2x2 center of mass)
		if fullRowLength > 2 && fullRows > 2 {
			remainingCells++
			partialRows = fullRows - 1
			fullRowLength--
		//Otherwise, if our row length is BIGGER than 3
		// (meaning if we pop off 2 and tack them onto the next row, it won't become a 3 row to a 2 parent)
		//AND if we have an idle row to put the cells in..
		} else if fullRowLength > 3 && fullRows < rows {
			//Just pop those 2 off the side
			remainingCells += 2
			fullRowLength--
		} else {
			fmt.Fprintf(outputFile, "Impossible\n")
			return
		}
	}
	
	var result []byte
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if i == 0 && j == 0 {
				result = append(result, 'c')
				continue
			} else if i < fullRows && j < fullRowLength || j == fullRowLength && i < partialRows || i == fullRows && j < remainingCells {
				result = append(result, '.')
			} else {
				result = append(result, '*')
			}
		}
		result = append(result, '\n')
	}
	
	
	fmt.Fprintf(outputFile, "%s", result)
}

func getSimpleSolution(outputFile *os.File, mines, rows, columns int) {
	//. empty
	//c clicked
	//* has a mine
	var result []byte
	for currentRow := 0; currentRow < rows; currentRow++ {
		for currentColumn := 0; currentColumn < columns; currentColumn++ {
			if currentRow == 0 && currentColumn == 0 {
				result = append(result, 'c')				
			} else {
				result = append(result, '*')
			}
		}
		result = append(result, '\n')
	}
	fmt.Fprintf(outputFile, "%s", result)
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