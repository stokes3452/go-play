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

func solveMaze(mazeSize int, maze [102][102]byte, interestingStuff []int) (result []byte) {
	benderY := interestingStuff[0]
	benderX :=interestingStuff[1]
	fryY := interestingStuff[2]
	fryX := interestingStuff[3]
	
	//0 - N
	//1 - E
	//2 - S
	//3 - W
	dir := 0
	if benderX == 1 && benderY == 1 {
		dir = 1
	} else if benderX == 1 {
		dir = 2
	} else if benderY == 1 {
		dir = 0
	} else {
		dir = 3
	}
	
	steps := 0
	//define some conditions here?
	//If we scan all 4 sides at the start, OR if we get back to our starting point, short-circuit... OR if we go too far :/
	for (steps < 4 || benderY != interestingStuff[0] || benderX != interestingStuff[1]) && len(result) < 10000 {
		steps++
		if dir >= 4 {
			dir -= 4
		} else if dir < 0 {
			dir += 4
		}
		Debug("Trying direction %d from %d,%d\n", dir, benderX, benderY)
		if benderX == fryX && benderY == fryY {
			Debug("Robo-mance achieved\n")
			return result
		}
		if dir == 0 {
			if maze[benderY-1][benderX] == '.' {
				benderY--
				result = append(result, 'N')
				dir--
			} else {
				dir++
			}
		} else if dir == 2 {
			if maze[benderY+1][benderX] == '.' {
				benderY++
				result = append(result, 'S')
				dir--
			} else {
				dir++
			}
		} else if dir == 3 {
			if maze[benderY][benderX-1] == '.' {
				benderX--
				result = append(result, 'W')
				dir--
			} else {
				dir++
			}
		} else if dir == 1 {
			if maze[benderY][benderX+1] == '.' {
				benderX++
				result = append(result, 'E')
				dir--
			} else {
				dir++
			}
		}
	}
	return []byte{}
}

func main() {
	testsRemaining := readIntLine()

	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		mazeSize := readIntLine()
		var maze [102][102]byte
		for i := 0; i < mazeSize; i++ {
			line := readLine()
			copy(maze[i+1][1:], line)
		}
		
		interestingStuff := getIntsFromLine()
		solution := solveMaze(mazeSize, maze, interestingStuff)
		
		if len(solution) > 0 {
			fmt.Fprintf(outputFile, "Case #%d: %d\n%s\n", testNumber, len(solution), solution)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: Edison ran out of energy.\n", testNumber)
		}
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