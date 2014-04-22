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
	var grid [100][100]room
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		dimensions := getIntsFromLine()
		height := dimensions[0]
		width := dimensions[1]
		Debug("Board size is %dx%d\n", height, width)
		remaining := getIntsFromLine()
		//The problem says X,Y X,Y, but the first example works either way
		//The problem set, however, only accepts this if you assume Y, X..
		entranceX := remaining[1]
		entranceY := remaining[0]
		exitX := remaining[3]
		exitY := remaining[2]
		//setup
		for yLoc := 0; yLoc < height; yLoc++ {
			yValues := getIntsFromLine()
			for xLoc := 0; xLoc < width; xLoc++ {
				grid[yLoc][xLoc].value = yValues[xLoc]
				grid[yLoc][xLoc].steps = -1
				grid[yLoc][xLoc].totalValue = -1
			}
		}
		Debug("Going from %d,%d to %d,%d\n", entranceY, entranceX, exitY, exitX)
		if entranceX >= width || exitX >= width || entranceY >= height || exitY >= height {
			fmt.Fprintf(outputFile, "Case #%d: Mission Impossible.\n", testNumber)
			continue

		}
		grid[entranceY][entranceX].steps = 0
		grid[entranceY][entranceX].totalValue = grid[entranceY][entranceX].value
		//now work on that whole solve nonsense
		bfsChannel := make(chan int, width*height)
		bfsChannel <- entranceY*100 + entranceX
		currentX := 0
		currentY := 0
		var currentRoom room
		var nextRoom *room
		var currentValue int
		ever := true
		for ever {
			select {
			case currentValue = <-bfsChannel:
				currentY = currentValue / 100
				currentX = currentValue % 100
				if currentX == exitX && currentY == exitY {
					ever = false
					break
				}
				currentRoom = grid[currentY][currentX]
				Debug("Current room corods is %d,%d\n", currentY, currentX)
				Debug("Current room value and steps is %d / %d\n", currentRoom.totalValue, currentRoom.steps)
				nextX, nextY := getNextRooms(height, width, currentX, currentY, grid)
				for nextIndex, _ := range nextX {
					Debug("Examining room %d,%d\n", nextY[nextIndex], nextX[nextIndex])
					nextRoom = &grid[nextY[nextIndex]][nextX[nextIndex]]
					if nextRoom.value >= 0 {
						if nextRoom.steps == -1 {
							nextRoom.steps = currentRoom.steps + 1
							Debug("Set next steps to %d\n", nextRoom.steps)
							bfsChannel <- nextY[nextIndex]*100 + nextX[nextIndex]
						}
						if nextRoom.steps == currentRoom.steps+1 {
							if nextRoom.value+currentRoom.totalValue > nextRoom.totalValue {
								nextRoom.totalValue = nextRoom.value + currentRoom.totalValue
							}
						}
					}
				}
				break
			default:
				ever = false
				break
			}
		}

		if grid[exitY][exitX].totalValue >= 0 {
			fmt.Fprintf(outputFile, "Case #%d: %d\n", testNumber, grid[exitY][exitX].totalValue)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: Mission Impossible.\n", testNumber)
		}
	}
}

func getNextRooms(height, width, currentX, currentY int, grid [100][100]room) (x, y []int) {
	if currentX > 0 {
		x = append(x, currentX-1)
		y = append(y, currentY)
	}
	if currentY > 0 {
		x = append(x, currentX)
		y = append(y, currentY-1)
	}
	if currentX < width-1 {
		x = append(x, currentX+1)
		y = append(y, currentY)
	}
	if currentY < height-1 {
		x = append(x, currentX)
		y = append(y, currentY+1)
	}
	return
}

type room struct {
	steps      int
	value      int
	totalValue int
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
