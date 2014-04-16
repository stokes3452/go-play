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
		problemScope := values[0]
		if problemScope == 1 {
			//we're dealing with a bitmap here, over the decision to move right or not
			//9th element = 1001, meaning right (1/1), left(1/2), left(1,3), right(4/3)
			index := values[1]
			var p uint64
			var q uint64
			p = 0
			q = 0
			for bit := 63; bit >= 0; bit-- {
				//If our bit is on, we're clear
				if (index & (1<<uint64(bit))) != 0 {
					//p, or start
					if p == 0 {
						p = 1
						q = 1
					} else {
						p = p + q
					}
					Debug("Shifted to %d from %d!\n", index & 1<<uint64(bit), index)
				} else {
					//q it is
					q = p + q
				}
			}
			fmt.Fprintf(outputFile, "Case #%d: %d %d\n", testNumber, p, q)
		} else if problemScope == 2 {
			//find the element's index
			p := values[1]
			q := values[2]
			var nodeSum uint64
			nodeSum = 0
			var nodeMultiplier uint64
			nodeMultiplier = 1
			for p != 1 || q != 1 {
				if q >  p {
					q = q - p
				} else {
					nodeSum += nodeMultiplier
					p = p - q
				}
				nodeMultiplier *= 2
			}
			nodeSum += nodeMultiplier
			fmt.Fprintf(outputFile, "Case #%d: %d\n", testNumber, nodeSum)
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

func getIntsFromLine() (lineValues []uint64) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.ParseUint(string(part), 10, 64)
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