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
	var board [4][4]byte
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		for rowNum := 0; rowNum < 4; rowNum++ {
			row := readLine()
			copy(board[rowNum][0:4], row)
		}
		if checkWinner(board, 'X') {
			fmt.Fprintf(outputFile, "Case #%d: X won\n", testNumber)
		} else if checkWinner(board, 'O') {
			fmt.Fprintf(outputFile, "Case #%d: O won\n", testNumber)
		} else {
			gameOver := true
			for _, row := range board {
				for _, column := range row {
					if column == '.' {
						if gameOver {
							fmt.Fprintf(outputFile, "Case #%d: Game has not completed\n", testNumber)
							gameOver = false
						}
					}
				}
			}
			if gameOver {
				fmt.Fprintf(outputFile, "Case #%d: Draw\n", testNumber)
			}
		}
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
		reader.ReadLine()
	}
}

func checkWinner(board [4][4]byte, tile byte) bool {
	for i := 0; i < 4; i++ {
		winningHorizontal := true
		winningVertical := true
		for j := 0; j < 4; j++ {
			if board[i][j] != tile && board[i][j] != 'T' {
				winningHorizontal = false
				Debug("No horizontal win because of %d,%d\n", i, j)
			}
			if board[j][i] != tile && board[j][i] != 'T' {
				winningVertical = false
				Debug("No vertical win because of %d,%d\n", j, i)
			}
		}
		if winningHorizontal || winningVertical {
			Debug("Natural win!\n")
			return true
		}
	}
	winningLeftDiagonal := true
	winningRightDiagonal := true
	for i := 0; i < 4; i++ {
		if board[i][i] != tile && board[i][i] != 'T' {
			winningRightDiagonal = false
		}
		if board[i][3-i] != tile && board[i][3-i] != 'T' {
			winningLeftDiagonal = false
		}
	}
	Debug("Diagonal win\n")
	return winningRightDiagonal || winningLeftDiagonal
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
