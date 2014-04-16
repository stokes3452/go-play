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
		numPairs := readIntLine()
		Debug("We're working with %d troublesome pairs\n", numPairs)
		//map of name into their bucket
		//buckets have left or right
		//IF neither name exists yet, make a new bucket
		//IF one does, add opponent to his bucket
		//IF opponent exists in another, merge!
		pairMap := make(map[string]pair)
		success := true
		for i := 0; i < numPairs; i++ {
			line := readLine()
			lineValues := bytes.Split(line, []byte{' '})
			firstName := string(lineValues[0])
			secondName := string(lineValues[1])
			firstPair, firstPairExists := pairMap[firstName]
			if !firstPairExists {
				firstPair = pair{}
				firstPair[firstName] = true
				pairMap[firstName] = firstPair
			}
			//Now, find out about our second pair...
			secondPair, secondPairExists := pairMap[secondName]
			if !secondPairExists {
				secondPair = pair{}
				secondPair[secondName] = true
				pairMap[secondName] = secondPair
			}
			
			//We now have a first pair AND a second pair
			//Zip our second pair into the first one
			//Start out by seeing how to combine them..
			combineOrder := (firstPair[firstName] == secondPair[secondName])
			for henchmanName, henchmanSide := range secondPair {
				//If at any point in time we reach a case where something exists on the wrong side of the fence, DANGER
				//Golang doesn't have a local xor, so we have to get a bit fugly with alternative != syntax
				newSide := (henchmanSide != combineOrder)
				oldSide, exists := firstPair[henchmanName]
				if exists {
					Debug("EXISTS!\n")
					if oldSide != newSide {
						Debug("Hench name is %s, Order of operation is %t, newSide is %t, oldSide is %t\n", henchmanName, combineOrder, newSide, oldSide)
						success = false
						break;
					}
				}
				firstPair[henchmanName] = newSide
				pairMap[henchmanName] = firstPair
			}
		}
		if success {
			fmt.Fprintf(outputFile, "Case #%d: Yes\n", testNumber)
		} else {
			fmt.Fprintf(outputFile, "Case #%d: No\n", testNumber)
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