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
		gameSize := readIntLine()
		var gameBoard [100][100]byte
		blueCount := 0
		redCount := 0
		for i := 0; i < gameSize; i++ {
			gameLine := readLine()
			copy(gameBoard[i][:], gameLine)
			for _, value := range gameLine {
				if value == 'B' {
					blueCount++

				} else if value == 'R' {
					redCount++
				}
			}
		}

		if blueCount+1 < redCount || redCount+1 < blueCount {
			fmt.Fprintf(outputFile, "Case #%d: Impossible\n", testNumber)
		} else {
			redWins := countWins(gameBoard, gameSize, 'R')
			gameBoard = invertBoard(gameBoard, gameSize)
			blueWins := countWins(gameBoard, gameSize, 'B')
			Debug("Win count is at %d, %d\n", redWins, blueWins)
			if blueWins >= 1 {
				if redCount > blueCount || blueWins > 1 {
					fmt.Fprintf(outputFile, "Case #%d: Impossible\n", testNumber)
				} else {
					fmt.Fprintf(outputFile, "Case #%d: Blue wins\n", testNumber)
				}
			} else if redWins >= 1 {
				if blueCount > redCount || redWins > 1 {
					fmt.Fprintf(outputFile, "Case #%d: Impossible\n", testNumber)
				} else {
					fmt.Fprintf(outputFile, "Case #%d: Red wins\n", testNumber)
				}
			} else {
				fmt.Fprintf(outputFile, "Case #%d: Nobody wins\n", testNumber)
			}
		}
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func getNextHex(x, y, gameSize, direction int) (int, int) {
	switch direction {
	case 0:
		x++
		break
	case 1:
		y++
		break
	case 2:
		y++
		x--
		break
	case 3:
		x--
		break
	case 4:
		y--
	case 5:
		y--
		x++
	}
	return x, y
}

func countWins(gameBoard [100][100]byte, gameSize int, tileCheck byte) (wins int) {
	if gameSize == 1 && gameBoard[0][0] == tileCheck {
		return 1
	}
	var locationStack [10000]int
	var stackLocation int
	//direction will be 0 through 5, 0 being right, and incrementing clockwise
	var directionStack [10000]int
	var direction int
	for x := 0; x < gameSize; x++ {
		stackLocation = 0
		locationStack[stackLocation] = x
		directionStack[stackLocation] = 2
		stackLocation++
		//If there's anything left, examine it
		for stackLocation > 0 {
			stackLocation--
			direction = directionStack[stackLocation]
			currentLocation := locationStack[stackLocation]
			currentX := currentLocation % gameSize
			currentY := currentLocation / gameSize
			//If our tile is lit up, check the tiles around it
			if gameBoard[currentY][currentX] == tileCheck {
				Debug("Nulling out %d, %d\n", currentY, currentX)
				gameBoard[currentY][currentX] = 0
				if currentY == gameSize-1 {
					wins++
					stackLocation = 0
				} else {
					gameBoard[currentY][currentX] = 0
					for change := -1; change < 2; change++ {
						currentDirection := int(direction) + change
						if currentDirection == -1 {
							currentDirection = 5
						} else if currentDirection == 6 {
							currentDirection = 0
						}
						//Find the next x and y
						nextX, nextY := getNextHex(currentX, currentY, gameSize, currentDirection)
						if nextX >= 0 && nextX < gameSize && nextY >= 0 && nextY < gameSize {
							directionStack[stackLocation] = currentDirection
							locationStack[stackLocation] = nextY*gameSize + nextX
							stackLocation++
						}
					}
				}
			}
		}
	}
	return
}

func invertBoard(gameBoard [100][100]byte, gameSize int) (newBoard [100][100]byte) {
	for i := 0; i < gameSize; i++ {
		for j := 0; j < gameSize; j++ {
			newBoard[i][j] = gameBoard[j][i]
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
		panic(fmt.Sprintf("Err reading int value from line: %d", line))
	}
	return
}
