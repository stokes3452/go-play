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
	Debug("Debug is %t\r\n", *debug)
	Debug("File in is: %s\r\n", *fileIn)
	Debug("File out is: %s\r\n", *fileOut)
	
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
	line := readLineOrPanic()
	
	testsRemaining, err := strconv.Atoi(string(line))
	if err != nil {
		panic("Err on first Atoi for test count")
	}
	
	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\r\n", testNumber)
		//Number of tests doesn't really matter.  We'll panic if we don't have the same number of blocks, regardless
		readLineOrPanic()

		//Now, get to work..
		myBlocks := getFloatsFromLine()
		opponentBlocks := getFloatsFromLine()
		if len(myBlocks) !=  len(opponentBlocks) {
			panic(fmt.Sprintf("Someone's cheating.. I have %d blocks and my opponent has %d blocks", myBlocks, opponentBlocks))
		}
		
		sort.Sort(myBlocks)
		sort.Sort(opponentBlocks)
		cheatingWins := 0
		cheatingLosses := 0
		realWins := 0
		for i := 0; i < len(myBlocks); i++ {
			//If we're cheating, we want to maximize our losses by popping off his biggest tiles with our worst ones
			if myBlocks[i] < opponentBlocks[i - cheatingLosses] {
				cheatingLosses++
			} else {
				cheatingWins++
			}
			//If we're not cheating, our opponent will win by the smallest margin possible, effectively shifting 1 position each time he has to lose
			if myBlocks[i - realWins] > opponentBlocks[i] {
				realWins++
			}
		}
		Debug("I just cheated, and won %d times!\r\n", cheatingWins)
		Debug("If I had not cheated, I would have only won %d times\r\n", realWins)
		fmt.Fprintf(outputFile, "Case #%d: %d %d\r\n", testNumber, cheatingWins, realWins)
		fmt.Printf("Processed %d of %d tests\r\n", testNumber, testsRemaining)
	}
}

func getFloatsFromLine() (lineValues sort.Float64Slice) {
	line := readLineOrPanic()
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
	line := readLineOrPanic()
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

func readLineOrPanic() (line []byte) {
	line, prefix, err := reader.ReadLine()
	if err != nil {
		panic("Error reading line for test")
	}
	if prefix {
		panic("Did not complete read.  Increase buffer size and try try again")
	}
	return
}